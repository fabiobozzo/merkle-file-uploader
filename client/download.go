package client

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file by index, from the server, and verify its integrity",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download...")
	},
}
