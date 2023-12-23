package storage

import "sync"

type InMemoryStorage struct {
	mu    sync.RWMutex
	seq   int
	files map[int]StoredFile
}

func (s *InMemoryStorage) StoreFile(file StoredFile) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	file.Index = s.seq
	s.files[s.seq] = file

	return s.seq, nil
}

func (s *InMemoryStorage) RetrieveAllFiles() (storedFiles []StoredFile, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := 1; i <= s.seq; i++ {
		storedFiles = append(storedFiles, s.files[i])
	}

	return
}

func (s *InMemoryStorage) RetrieveFileByIndex(i int) (storedFile StoredFile, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storedFile, found := s.files[i]
	if !found {
		err = ErrStoredFileNotFound
	}

	return
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		files: make(map[int]StoredFile),
	}
}
