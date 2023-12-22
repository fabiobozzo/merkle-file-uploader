package download

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	ErrFailedDownload = errors.New("failed to download file")
)

type HttpDownloader struct {
	Client  *http.Client
	BaseURL string
}

func NewHttpDownloader(httpClient *http.Client, baseURL string) *HttpDownloader {
	return &HttpDownloader{
		Client:  httpClient,
		BaseURL: baseURL,
	}
}

func (h *HttpDownloader) DownloadFileAt(index int, destination *os.File) (err error) {
	response, err := http.Get(fmt.Sprintf("%s/download/%d", h.BaseURL, index))
	if err != nil {
		err = fmt.Errorf("%w: error sending GET request: %s", ErrFailedDownload, err)

		return
	}
	defer func() { _ = response.Body.Close() }()

	if _, err = io.Copy(destination, response.Body); err != nil {
		err = fmt.Errorf("%w: error reading response body: %s", ErrFailedDownload, err)
	}

	return
}
