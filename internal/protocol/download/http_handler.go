package download

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"merkle-file-uploader/internal/merkle"
	"merkle-file-uploader/internal/protocol"
	"merkle-file-uploader/internal/storage"
	"merkle-file-uploader/internal/utils"
)

func NewDownloadHandler(repository storage.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.HttpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		index, err := indexFromRequest(r)
		if err != nil {
			utils.HttpError(w, http.StatusBadRequest, err)

			return
		}

		fileContent, err := repository.RetrieveFileByIndex(index)
		if err == storage.ErrStoredFileNotFound {
			utils.HttpError(w, http.StatusNotFound, fmt.Errorf("{index} not found: %d", index))

			return
		}
		if err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)

			return
		}

		_, err = w.Write(fileContent.Content)

		return
	}
}

func NewProofHandler(repository storage.Repository, hashFn merkle.HashFn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.HttpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		index, err := indexFromRequest(r)
		if err != nil {
			utils.HttpError(w, http.StatusBadRequest, err)

			return
		}

		storedFiles, err := repository.RetrieveAllFiles()
		if err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)

			return
		}

		var blocks []string
		var blockToProve string
		for _, f := range storedFiles {
			fileContent := string(f.Content)
			blocks = append(blocks, fileContent)
			if f.Index == index {
				blockToProve = fileContent
			}
		}

		merkleTree, err := merkle.NewTree(blocks, hashFn)
		if err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)

			return
		}

		merkleProof := merkleTree.ProofForBlock(blockToProve)
		if err = utils.HttpOkJson(w, protocol.MerkleProofResponse{MerkleProof: merkleProof}); err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)
		}

		log.Println("merkle root:", merkleTree.Root.Data)
		log.Println("merkle proof:", merkleProof)

		return
	}
}

func indexFromRequest(r *http.Request) (index int, err error) {
	vars := mux.Vars(r)
	indexParam, isIndexSet := vars["index"]
	if !isIndexSet {
		err = errors.New("{index} path param is not passed in")

		return
	}

	index, err = strconv.Atoi(indexParam)
	if err != nil {
		err = fmt.Errorf("{index} path param must be numeric: %s", err)
	}

	return
}
