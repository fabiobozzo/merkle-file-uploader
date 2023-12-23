package storage

import "errors"

var (
	ErrStoredFileNotFound = errors.New("the file is not found in the storage")
)

type StoredFile struct {
	Index   int
	Name    string
	Content []byte
}

type Repository interface {
	StoreFile(file StoredFile) (int, error)
	RetrieveAllFiles() ([]StoredFile, error)
	RetrieveFileByIndex(i int) (StoredFile, error)
}
