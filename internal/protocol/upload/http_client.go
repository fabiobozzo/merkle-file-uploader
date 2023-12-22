package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"merkle-file-uploader/internal/protocol"
)

var (
	ErrFailedUpload = errors.New("failed to upload files")
)

type HttpUploader struct {
	Client  *http.Client
	BaseURL string
}

func NewHttpUploader(httpClient *http.Client, baseURL string) *HttpUploader {
	return &HttpUploader{
		Client:  httpClient,
		BaseURL: baseURL,
	}
}

func (h *HttpUploader) UploadFilesFrom(filePaths []string) (uploadedFiles []protocol.UploadedFile, err error) {
	requestBody, formDataContentType, err := prepareRequestBody(filePaths)
	if err != nil {
		err = fmt.Errorf("%w: error preparing POST request body: %s", ErrFailedUpload, err)

		return
	}

	response, err := http.Post(fmt.Sprintf("%s/upload", h.BaseURL), formDataContentType, &requestBody)
	if err != nil {
		err = fmt.Errorf("%w: error sending POST request: %s", ErrFailedUpload, err)

		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("%w: unexpected http status: %s", ErrFailedUpload, response.Status)

		return
	}

	var decodedResponse protocol.UploadedFilesResponse
	if err = json.NewDecoder(response.Body).Decode(&decodedResponse); err != nil {
		err = fmt.Errorf("%w: error decoding json response: %s", ErrFailedUpload, err)

		return
	}

	defer func() { _ = response.Body.Close() }()

	return decodedResponse.UploadedFiles, nil
}

func prepareRequestBody(filePaths []string) (requestBody bytes.Buffer, formDataContentType string, err error) {
	multipartWriter := multipart.NewWriter(&requestBody)

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
