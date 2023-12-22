package download

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"merkle-file-uploader/internal/storage"
	"merkle-file-uploader/internal/utils"
)

func NewDownloadHandler(repository storage.Repository) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.HttpError(w, http.StatusMethodNotAllowed, errors.New(r.Method))

			return
		}

		vars := mux.Vars(r)
		indexParam, isIndexSet := vars["index"]
		if !isIndexSet {
			utils.HttpError(w, http.StatusBadRequest, errors.New("{index} path param is not passed in"))

			return
		}

		index, err := strconv.Atoi(indexParam)
		if err != nil {
			utils.HttpError(w, http.StatusBadRequest, fmt.Errorf("{index} path param must be numeric: %s", err))

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
