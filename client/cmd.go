package client

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"merkle-file-uploader/internal/utils"
)

const (
	defaultServerURL          = "http://localhost:8080"
	defaultMerkleRootFilename = ".merkleroot"
)

var hashFn = utils.Sha256

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "The mfu client can upload & download files and verify their integrity",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Cmd.AddCommand(uploadCmd)
	Cmd.AddCommand(downloadCmd)
}

func getServerURL() (serverUrl string) {
	serverUrl = os.Getenv("SERVER_URL")
	if len(serverUrl) == 0 {
		serverUrl = defaultServerURL
	}

	return
}

func getMerkleRootFilename() (merkleRootFilename string) {
	merkleRootFilename = os.Getenv("MERKLE_ROOT_FILENAME")
	if len(merkleRootFilename) == 0 {
		merkleRootFilename = defaultMerkleRootFilename
	}

	return
}
