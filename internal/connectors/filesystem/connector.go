package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/drekunov/gc/internal/app"
)

const (
	connectorName = "FileSystem"

	attrName  = "Name"
	attrSize  = "Size"
	attrDate  = "Date"
	attrTime  = "Time"
	attrIsDir = "IsDir"
	attrKind  = "Kind"

	dirSizeValue = "<DIR>"

	dateLayout = "2006-01-02"
	timeLayout = "15:04:05"

	columnWidthName = 11
	columnWidthSize = 6
	columnWidthDate = 10
	columnWidthTime = 8
)

// kindByExtension maps a lowercase file extension, including its leading dot,
// to the category kind it belongs to.
var kindByExtension = buildKindByExtension(map[string][]string{
	app.FileKindImage: {
		"png", "jpg", "jpeg", "gif", "bmp", "svg", "webp", "ico", "tif", "tiff",
	},
	app.FileKindArchive: {
		"zip", "tar", "gz", "bz2", "xz", "zst", "7z", "rar", "tgz",
	},
	app.FileKindSource: {
		"go", "c", "h", "cc", "cpp", "hpp", "rs", "py", "js", "mjs", "cjs",
		"ts", "tsx", "jsx", "java", "kt", "kts", "rb", "php", "sh", "bash",
		"zsh", "lua", "pl", "r", "swift", "cs", "scala", "hs", "sql", "vim",
	},
	app.FileKindConfig: {
		"json", "yaml", "yml", "toml", "ini", "cfg", "conf", "env",
	},
})

// buildKindByExtension inverts a category-to-extensions grouping into an
// extension-to-kind lookup, adding the leading dot to every extension.
func buildKindByExtension(groups map[string][]string) map[string]string {
	index := make(map[string]string)

	for kind, extensions := range groups {
		for _, extension := range extensions {
			index["."+extension] = kind
		}
	}

	return index
}

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
		{AttrName: attrKind, AttrValue: attrKind, Hidden: true},
	})

	for _, entry := range entries {
		attributeList = append(attributeList, entryAttributes(path, entry))
	}

	return attributeList, nil
}

// entryAttributes builds the listing row for one entry: its name, a
// human-readable size (or <DIR> for a directory), and the modification date and
// time. The directory flag and file kind are carried as hidden attributes.
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
		{AttrName: attrKind, AttrValue: fileKind(path, entry), Hidden: true},
	}
}

// fileKind classifies an entry as a directory, a symlink, an executable, or an
// extension-based category, choosing the first that matches, and returns "" when
// none does. A symlink that resolves to a directory is a directory, matching
// the navigation behavior.
func fileKind(path string, entry os.DirEntry) string {
	if isDirEntry(path, entry) {
		return app.FileKindDirectory
	}

	if entry.Type()&os.ModeSymlink != 0 {
		return app.FileKindSymlink
	}

	if isExecutable(entry) {
		return app.FileKindExecutable
	}

	return kindFromName(entry.Name())
}

// isExecutable reports whether the entry is a regular file with an execute
// permission bit. A non-regular file is not an executable.
func isExecutable(entry os.DirEntry) bool {
	info, err := entry.Info()
	if err != nil {
		return false
	}

	mode := info.Mode()

	return mode.IsRegular() && mode.Perm()&0o111 != 0
}

// kindFromName maps a file name's extension to a category kind, matching
// case-insensitively and returning "" for an unknown extension.
func kindFromName(name string) string {
	return kindByExtension[strings.ToLower(filepath.Ext(name))]
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
