package storage

import (
	"context"
	"errors"
)

var (
	ErrStoredFileNotFound = errors.New("the file is not found in the storage")
)

type StoredFile struct {
	Index   int
	Name    string
	Content []byte
}

type Repository interface {
	StoreFile(ctx context.Context, file StoredFile) (int, error)
	RetrieveAllFiles(ctx context.Context) ([]StoredFile, error)
	RetrieveFileByIndex(ctx context.Context, i int) (StoredFile, error)
	DeleteAllFiles(ctx context.Context) error
}
