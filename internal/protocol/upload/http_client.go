package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"merkle-file-uploader/internal/merkle"
	"merkle-file-uploader/internal/protocol"
	"merkle-file-uploader/internal/utils"
)

var (
	ErrFailedUpload = errors.New("failed to upload files")
)

type HttpUploader struct {
	client  *http.Client
	baseURL string
	hashFn  merkle.HashFn
}

func NewHttpUploader(httpClient *http.Client, baseURL string, hashFn merkle.HashFn) *HttpUploader {
	return &HttpUploader{
		client:  httpClient,
		baseURL: baseURL,
		hashFn:  hashFn,
	}
}

func (h *HttpUploader) UploadFilesFrom(filePaths []string) (
	uploadedFiles []protocol.UploadedFile,
	merkleRoot string,
	err error,
) {
	requestBody, formDataContentType, err := utils.MultipartFormFromFiles(filePaths)
	if err != nil {
		err = fmt.Errorf("%w: error preparing POST request body: %s", ErrFailedUpload, err)

		return
	}

	response, err := http.Post(fmt.Sprintf("%s/upload", h.baseURL), formDataContentType, &requestBody)
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

	merkleRoot, err = h.computeMerkleRoot(filePaths)
	if err != nil {
		err = fmt.Errorf("%w: error computing merkle root: %s", ErrFailedUpload, err)

		return
	}

	return decodedResponse.UploadedFiles, merkleRoot, nil
}

func (h *HttpUploader) computeMerkleRoot(filePaths []string) (merkleRoot string, err error) {
	var blocks []string
	for _, f := range filePaths {
		var fileContent []byte
		fileContent, err = os.ReadFile(f)
		if err != nil {
			err = fmt.Errorf("%w: error reading file for hashing: %s", ErrFailedUpload, err)

			return
		}

		blocks = append(blocks, string(fileContent))
	}

	tree, err := merkle.NewTree(blocks, h.hashFn)
	if err != nil {
		return
	}

	return tree.Root.Data, nil
}
