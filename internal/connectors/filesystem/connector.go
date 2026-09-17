package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/drekunov/gc/internal/app"
)

const (
	connectorName = "FileSystem"

	attrName  = "Name"
	attrIsDir = "IsDir"
	attrType  = "Type"
)

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

func (f *FileSystem) ReadDir(path string) ([]app.AttributeList, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	if len(entries) == 0 {
		return nil, nil
	}

	attributeList := make([]app.AttributeList, 0, len(entries)+1)

	attributeList = append(attributeList, app.AttributeList{
		{AttrName: attrName, AttrValue: attrName},
		{AttrName: attrIsDir, AttrValue: attrIsDir},
		{AttrName: attrType, AttrValue: attrType},
	})

	for _, entry := range entries {
		attrs := make(app.AttributeList, 0, 3)

		attrs = append(attrs, app.Attribute{
			AttrName:  attrName,
			AttrValue: entry.Name(),
		})

		attrs = append(attrs, app.Attribute{
			AttrName:  attrIsDir,
			AttrValue: isDirEntry(path, entry),
		})

		attrs = append(attrs, app.Attribute{
			AttrName:  attrType,
			AttrValue: entry.Type(),
		})

		attributeList = append(attributeList, attrs)
	}

	return attributeList, nil
}

// isDirEntry reports whether the entry is a directory, following a symbolic
// link to its target. A link that does not resolve to a directory, or whose
// target cannot be stat'd, is not a directory.
func isDirEntry(dir string, entry os.DirEntry) bool {
	if entry.IsDir() {
		return true
	}

	if entry.Type()&os.ModeSymlink == 0 {
		return false
	}

	info, err := os.Stat(filepath.Join(dir, entry.Name()))
	if err != nil {
		return false
	}

	return info.IsDir()
}
