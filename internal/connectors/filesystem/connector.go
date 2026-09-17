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
	attrSize  = "Size"
	attrDate  = "Date"
	attrTime  = "Time"
	attrIsDir = "IsDir"

	dirSizeValue = "<DIR>"

	dateLayout = "2006-01-02"
	timeLayout = "15:04:05"

	columnWidthName = 11
	columnWidthSize = 6
	columnWidthDate = 10
	columnWidthTime = 8
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
		{AttrName: attrName, AttrValue: attrName, Width: columnWidthName, Flex: true},
		{AttrName: attrSize, AttrValue: attrSize, Width: columnWidthSize},
		{AttrName: attrDate, AttrValue: attrDate, Width: columnWidthDate},
		{AttrName: attrTime, AttrValue: attrTime, Width: columnWidthTime},
		{AttrName: attrIsDir, AttrValue: attrIsDir, Hidden: true},
	})

	for _, entry := range entries {
		attributeList = append(attributeList, entryAttributes(path, entry))
	}

	return attributeList, nil
}

// entryAttributes builds the listing row for one entry: its name, a
// human-readable size (or <DIR> for a directory), and the modification date and
// time. The directory flag is carried as a hidden attribute.
func entryAttributes(path string, entry os.DirEntry) app.AttributeList {
	isDir := isDirEntry(path, entry)

	var size any = ""
	if isDir {
		size = dirSizeValue
	}

	date := ""
	clock := ""

	info, err := entry.Info()
	if err == nil {
		if !isDir {
			size = app.Size(info.Size())
		}

		modTime := info.ModTime()
		date = modTime.Format(dateLayout)
		clock = modTime.Format(timeLayout)
	}

	return app.AttributeList{
		{AttrName: attrName, AttrValue: entry.Name()},
		{AttrName: attrSize, AttrValue: size},
		{AttrName: attrDate, AttrValue: date},
		{AttrName: attrTime, AttrValue: clock},
		{AttrName: attrIsDir, AttrValue: isDir, Hidden: true},
	}
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
