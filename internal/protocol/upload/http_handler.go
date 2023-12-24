package upload

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"merkle-file-uploader/internal/merkle"
	"merkle-file-uploader/internal/protocol"
	"merkle-file-uploader/internal/storage"
	"merkle-file-uploader/internal/utils"
)

func NewUploadHandler(repository storage.Repository, hashFn merkle.HashFn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.HttpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		// limit maxMultipartMemory
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utils.HttpError(w, http.StatusBadRequest, fmt.Errorf("unable to parse multipart form: %s", err))

			return
		}

		if err := repository.DeleteAllFiles(r.Context()); err != nil {
			utils.HttpError(w, http.StatusInternalServerError, fmt.Errorf("error while resetting storage: %s", err))

			return
		}

		var uploadedFiles []protocol.UploadedFile
		var blocks []string

		files := r.MultipartForm.File["files"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				utils.HttpError(w, http.StatusBadRequest, fmt.Errorf("unable to open file: %s", err))

				return
			}
			_ = file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				utils.HttpError(w, http.StatusBadRequest, fmt.Errorf("unable to read file: %s", err))

				return
			}

			i, err := repository.StoreFile(r.Context(), storage.StoredFile{
				Name:    fileHeader.Filename,
				Content: data,
			})
			if err != nil {
				utils.HttpError(w, http.StatusInternalServerError, err)

				return
			}

			uploadedFiles = append(uploadedFiles, protocol.UploadedFile{
				Name:  fileHeader.Filename,
				Index: i,
			})

			blocks = append(blocks, string(data))
		}

		merkleTree, err := merkle.NewTree(blocks, hashFn)
		if err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)

			return
		}

		if err = repository.StoreTree(r.Context(), merkleTree); err != nil {
			utils.HttpError(w, http.StatusInternalServerError, fmt.Errorf("unable to store the merkle tree: %s", err))

			return
		}

		if err := utils.HttpOkJson(w, protocol.UploadedFilesResponse{UploadedFiles: uploadedFiles}); err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)

			return
		}
	}
}
