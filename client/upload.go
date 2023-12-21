package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a set of files, or an entire folder, to the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upload...")
	},
}
