package utils

import (
	"os"
	"path/filepath"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), nil
}

func ListFilesInDirectory(directoryPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude directories
		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
