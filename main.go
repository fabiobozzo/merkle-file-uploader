package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"merkle-file-uploader/client"
	"merkle-file-uploader/server"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mfu",
		Short: "MFU is a tool for verifiable file uploads and downloads",
		Long: `A CLI for managing both MFU client and server,
			and verifying downloaded files with the help of Merkle proofs.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}
		},
	}

	rootCmd.AddCommand(client.Cmd)
	rootCmd.AddCommand(server.Cmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
