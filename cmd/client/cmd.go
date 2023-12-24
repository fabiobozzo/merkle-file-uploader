package client

import (
	"log"

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
