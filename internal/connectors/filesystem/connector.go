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
		attrs := make(app.AttributeList, 0, 4)

		attrs = append(attrs, app.Attribute{
			AttrName:  "path",
			AttrValue: path,
		})

		attrs = append(attrs, app.Attribute{
			AttrName:  "Name",
			AttrValue: entry.Name(),
		})

		attrs = append(attrs, app.Attribute{
			AttrName:  "IsDir",
			AttrValue: entry.IsDir(),
		})

		attrs = append(attrs, app.Attribute{
			AttrName:  "Type",
			AttrValue: entry.Type(),
		})

		attributeList = append(attributeList, attrs)
	}

	return attributeList, nil
}
