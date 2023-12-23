package client

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"merkle-file-uploader/internal/protocol/download"
)

type Downloader interface {
	DownloadFileAt(index int, destination *os.File) error
}

var _ Downloader = (*download.HttpDownloader)(nil)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file by index, from the server, and verify its integrity",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please enter one index of an uploaded file to download")

			return
		}

		index, err := strconv.Atoi(args[0])
		if err != nil || index < 1 {
			fmt.Println("The index must be a number starting from 1")

			return
		}

		rootHash, err := os.ReadFile(getMerkleRootFilename())
		if err != nil {
			fmt.Println("Merkle Root hash is missing or unreadable:", err)

			return
		}

		downloader := download.NewHttpDownloader(
			&http.Client{Timeout: time.Second * 30},
			getServerURL(),
			string(rootHash),
			hashFn,
		)

		if err := downloader.DownloadFileAt(index, os.Stdout); err != nil {
			fmt.Println(err)

			return
		}
	},
}
