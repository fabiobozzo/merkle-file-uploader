package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"merkle-file-uploader/internal/protocol/download"
	"merkle-file-uploader/internal/protocol/upload"
	"merkle-file-uploader/internal/storage"
	"merkle-file-uploader/internal/utils"
)

const defaultPort = 8080

var hashFn = utils.Sha256

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "The mfu server exposes a HTTP API for verifiable files upload & download",
	Run: func(cmd *cobra.Command, args []string) {
		inMemoryStorage := storage.NewInMemoryStorage()

		r := mux.NewRouter()
		r.HandleFunc("/upload", upload.NewUploadHandler(inMemoryStorage))
		r.HandleFunc("/download/{index}", download.NewDownloadHandler(inMemoryStorage))
		r.HandleFunc("/proof/{index}", download.NewProofHandler(inMemoryStorage, hashFn))

		log.Println("mfu server started on port", getPort())
		if err := http.ListenAndServe(fmt.Sprintf(":%d", getPort()), r); err != nil {
			log.Fatal(err)
		}
	},
}

func getPort() (port int) {
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = defaultPort
	}

	return
}
