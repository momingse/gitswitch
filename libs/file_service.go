package libs

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) GetCurrentPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func (f *FileService) GetParentFolderName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

func (f *FileService) CheckIfPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}
