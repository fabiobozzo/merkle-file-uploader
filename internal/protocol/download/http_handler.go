package download

import (
	"errors"
	"fmt"
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

		fileContent, err := repository.RetrieveFileByIndex(r.Context(), index)
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

		merkleTree, err := repository.RetrieveTree(r.Context())
		if err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)

			return
		}
		merkleTree.HashFn = hashFn

		fileByIndex, err := repository.RetrieveFileByIndex(r.Context(), index)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, storage.ErrStoredFileNotFound) {
				statusCode = http.StatusNotFound
			}

			utils.HttpError(w, statusCode, err)

			return
		}

		merkleProof := merkleTree.ProofForBlock(string(fileByIndex.Content))
		if err = utils.HttpOkJson(w, protocol.MerkleProofResponse{MerkleProof: merkleProof}); err != nil {
			utils.HttpError(w, http.StatusInternalServerError, err)
		}

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
