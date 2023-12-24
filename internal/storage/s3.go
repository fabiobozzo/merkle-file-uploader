package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(accessKeyId, secretAccessKey, endpoint, bucket string) (s3Storage *S3Storage, err error) {
	s3Storage = &S3Storage{
		bucket: bucket,
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: "us-east-1",
			}, nil
		})),
	)
	if err != nil {
		return
	}

	s3Storage.client = s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})

	return
}

func (s *S3Storage) StoreFile(ctx context.Context, file StoredFile) (i int, err error) {
	filesCount, err := s.countFiles(ctx)
	if err != nil {
		return
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%d", filesCount+1)),
		Body:   bytes.NewReader(file.Content),
	})

	return
}

func (s *S3Storage) RetrieveAllFiles(ctx context.Context) (storedFiles []StoredFile, err error) {
	filesCount, err := s.countFiles(ctx)
	if err != nil {
		return
	}

	for i := 1; i <= filesCount; i++ {
		var f StoredFile
		f, err = s.RetrieveFileByIndex(ctx, i)
		if err != nil {
			return
		}

		storedFiles = append(storedFiles, f)
	}

	return
}

func (s *S3Storage) RetrieveFileByIndex(ctx context.Context, i int) (storedFile StoredFile, err error) {
	resp, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%d", i)),
	})
	var nsk *types.NoSuchKey
	if errors.As(err, &nsk) {
		err = ErrStoredFileNotFound

		return
	}
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	storedFile.Index = i
	storedFile.Name = fmt.Sprintf("%d", i) // TODO: preserve original filename in `StoreFile`
	storedFile.Content, err = io.ReadAll(resp.Body)

	return
}

func (s *S3Storage) DeleteAllFiles(ctx context.Context) (err error) {
	objs, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	if err != nil {
		return
	}

	for _, object := range objs.Contents {
		_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    object.Key,
		})
		if err != nil {
			return
		}
	}

	return nil
}

func (s *S3Storage) countFiles(ctx context.Context) (count int, err error) {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})
	for paginator.HasMorePages() {
		var page *s3.ListObjectsV2Output
		page, err = paginator.NextPage(ctx)
		if err != nil {
			return
		}

		count += len(page.Contents)
	}

	return
}
