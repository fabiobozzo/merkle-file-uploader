package download

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"merkle-file-uploader/internal/merkle"
	"merkle-file-uploader/internal/protocol"
)

var (
	ErrFailedDownload     = errors.New("failed to download file")
	ErrFailedVerification = errors.New("the file integrity is compromised")
)

type HttpDownloader struct {
	client   *http.Client
	baseURL  string
	rootHash string
	hashFn   merkle.HashFn
}

func NewHttpDownloader(httpClient *http.Client, baseURL, rootHash string, hashFn merkle.HashFn) *HttpDownloader {
	return &HttpDownloader{
		client:   httpClient,
		baseURL:  baseURL,
		rootHash: rootHash,
		hashFn:   hashFn,
	}
}

func (h *HttpDownloader) DownloadFileAt(index int, destination *os.File) (err error) {
	downloadResponse, err := http.Get(fmt.Sprintf("%s/download/%d", h.baseURL, index))
	if err != nil {
		err = fmt.Errorf("%w: error sending GET /download request: %s", ErrFailedDownload, err)

		return
	}
	defer func() { _ = downloadResponse.Body.Close() }()

	if downloadResponse.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("%w: file not found at index %d", ErrFailedDownload, index)

		return
	}

	proofResponse, err := http.Get(fmt.Sprintf("%s/proof/%d", h.baseURL, index))
	if err != nil {
		err = fmt.Errorf("%w: error sending GET /proof request: %s", ErrFailedDownload, err)

		return
	}
	defer func() { _ = proofResponse.Body.Close() }()

	var merkleProof protocol.MerkleProofResponse
	if err = json.NewDecoder(proofResponse.Body).Decode(&merkleProof); err != nil {
		err = fmt.Errorf("%w: error decoding merkle proof response body: %s", ErrFailedDownload, err)

		return
	}

	fileContent, err := io.ReadAll(downloadResponse.Body)
	if err != nil {
		err = fmt.Errorf("%w: error reading download response body %s", ErrFailedDownload, err)

		return
	}

	if verified := merkle.VerifyProof(h.rootHash, string(fileContent), merkleProof.MerkleProof, h.hashFn); !verified {
		err = fmt.Errorf("%w: merkle root does not match: %s", ErrFailedDownload, h.rootHash)

		return
	}

	reader := io.NopCloser(bytes.NewReader(fileContent))
	if _, err = io.Copy(destination, reader); err != nil {
		err = fmt.Errorf("%w: error reading downloaded file: %s", ErrFailedDownload, err)
	}

	return
}
