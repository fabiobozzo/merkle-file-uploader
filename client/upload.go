package client

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"merkle-file-uploader/internal/protocol"
	"merkle-file-uploader/internal/protocol/upload"
	"merkle-file-uploader/internal/utils"
)

type Uploader interface {
	UploadFilesFrom(filePaths []string) ([]protocol.UploadedFile, string, error)
}

var _ Uploader = (*upload.HttpUploader)(nil)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a set of files, or an entire folder, to the server",
	Long:  "E.g. args: <file1> <file2> <file3> | args: <directory>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please enter a list of file paths or a directory path to upload")

			return
		}

		filePaths, err := argsToFilesToUpload(args)
		if err != nil {
			fmt.Println(err)

			return
		}

		uploader := upload.NewHttpUploader(&http.Client{Timeout: time.Second * 30}, getServerURL(), hashFn)
		uploadedFiles, merkleRoot, err := uploader.UploadFilesFrom(filePaths)
		if err != nil {
			fmt.Println(err)

			return
		}

		for _, f := range uploadedFiles {
			fmt.Printf("Uploaded file at index #%d: %s\n", f.Index, f.Name)
		}

		if err = os.WriteFile(getMerkleRootFilename(), []byte(merkleRoot), 0644); err != nil {
			fmt.Printf("Failed to store merkle root: %s\n", err)

			return
		}

		fmt.Println("Merkle Root hash:", merkleRoot)
	},
}

func argsToFilesToUpload(args []string) (filePaths []string, err error) {
	// Check whether the 1st arg is a directory path
	isDirectory, err := utils.IsDirectory(args[0])
	if err != nil {
		err = fmt.Errorf("error checking if %s is a directory: %v\n", args[0], err)

		return
	}

	if isDirectory {
		filePaths, err = utils.ListFilesInDirectory(args[0])
		if err != nil {
			err = fmt.Errorf("error listing files inside of %s: %v\n", args[0], err)

			return
		}
	} else {
		for _, arg := range args {
			if _, err := os.Stat(arg); err == nil {
				filePaths = append(filePaths, arg)
			}
		}
	}

	if len(filePaths) == 0 {
		err = errors.New("none of the files/dir specified can be found")
	}

	return
}
