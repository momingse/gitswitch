package libs

import (
	"os"
	"path/filepath"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) getCurrentPath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func (f *FileService) getParentFolderName(path string) (string, error) {
	return filepath.Base(filepath.Dir(path)), nil
}
