package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func HttpOkJson(w http.ResponseWriter, payload any) (err error) {
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(payload)
}

func HttpError(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("%s: %s\n", http.StatusText(statusCode), err)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func MultipartFormFromFiles(filePaths []string) (multipartForm bytes.Buffer, formDataContentType string, err error) {
	multipartWriter := multipart.NewWriter(&multipartForm)

	for _, fp := range filePaths {
		var file *os.File
		file, err = os.Open(fp)
		if err != nil {
			return
		}

		var filePart io.Writer
		filePart, err = multipartWriter.CreateFormFile("files", filepath.Base(fp))
		if err != nil {
			return
		}

		// copy the file content to the form file part
		if _, err = io.Copy(filePart, file); err != nil {
			return
		}

		if err = file.Close(); err != nil {
			return
		}
	}

	// Close the multipart writer to finish building the request body
	if err = multipartWriter.Close(); err != nil {
		return
	}

	formDataContentType = multipartWriter.FormDataContentType()

	return
}
