package client

import (
	"log"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "client",
	Short: "The MFU client can upload & download files and verify their integrity",
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
