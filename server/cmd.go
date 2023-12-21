package server

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "The MFU server exposes a HTTP API for verifiable files upload & download",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server...")
	},
}
