package filesystem

import (
	"errors"
	"fmt"
	"os"

	"github.com/drekunov/gc/internal/app"
)

const connectorName = "FileSystem"

type FileSystem struct{}

func New() *FileSystem {
	return &FileSystem{}
}

func (f *FileSystem) ReadFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path) //nolint: gosec
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return file, nil
}

func (f *FileSystem) Name() string {
	return connectorName
}

var ErrDirectoryEmpty = errors.New("directory is empty")

func (f *FileSystem) ReadDir(path string) ([]app.AttributeList, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	if len(entries) == 0 {
		return nil, ErrDirectoryEmpty
	}

	attributeList := make([]app.AttributeList, 0, len(entries))

	for _, entry := range entries {
		attrs := make(app.AttributeList)

		attrs["path"] = path
		attrs["Name"] = entry.Name()
		attrs["IsDir"] = entry.IsDir()
		attrs["Type"] = entry.Type()

		info, err := entry.Info()
		if err != nil {
			attrs["error"] = err
		} else {
			attrs["ModTime"] = info.ModTime()
			attrs["Size"] = info.Size()
		}

		attributeList = append(attributeList, attrs)
	}

	return attributeList, nil
}
