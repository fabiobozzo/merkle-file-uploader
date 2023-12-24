package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"merkle-file-uploader/internal/protocol/download"
	"merkle-file-uploader/internal/protocol/upload"
	"merkle-file-uploader/internal/storage"
	"merkle-file-uploader/internal/utils"
)

const (
	defaultPort               = 8080
	defaultAwsAccessKeyId     = "test"
	defaultAwsSecretAccessKey = "test"
	defaultAwsEndpoint        = "http://localstack:4566"
	defaultS3BucketName       = "mfu-202312"
)

var hashFn = utils.Sha256

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "The mfu server exposes a HTTP API for verifiable files upload & download",
	Run: func(cmd *cobra.Command, args []string) {
		//repository := storage.NewInMemoryStorage()
		repository, err := storage.NewS3Storage(
			utils.EnvStr("AWS_ACCESS_KEY_ID", defaultAwsAccessKeyId),
			utils.EnvStr("AWS_SECRET_ACCESS_KEY", defaultAwsSecretAccessKey),
			utils.EnvStr("AWS_ENDPOINT", defaultAwsEndpoint),
			utils.EnvStr("AWS_S3_BUCKET_NAME", defaultS3BucketName),
		)
		if err != nil {
			log.Fatal("error while connecting to S3:", err)

			return
		}

		r := mux.NewRouter()
		r.HandleFunc("/upload", upload.NewUploadHandler(repository))
		r.HandleFunc("/download/{index}", download.NewDownloadHandler(repository))
		r.HandleFunc("/proof/{index}", download.NewProofHandler(repository, hashFn))

		port := utils.EnvInt("PORT", defaultPort)
		log.Println("mfu server started on port", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
			log.Fatal(err)
		}
	},
}
