package filesystem_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/connectors/filesystem"
)

const (
	attrName  = "Name"
	attrSize  = "Size"
	attrDate  = "Date"
	attrTime  = "Time"
	attrIsDir = "IsDir"
)

var (
	datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timePattern = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}$`)
)

func TestReadDirHeaderContract(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(dir, "sub"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := filesystem.New().ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 3 {
		t.Fatalf("expected header + 2 entries, got %d rows", len(rows))
	}

	assertHeader(t, rows[0])
	assertFileRow(t, rows[1])
	assertDirRow(t, rows[2])
}

func assertHeader(t *testing.T, header app.AttributeList) {
	t.Helper()

	wantHeader := []struct {
		name   string
		width  int
		hidden bool
		flex   bool
	}{
		{name: attrName, width: 11, flex: true},
		{name: attrSize, width: 6},
		{name: attrDate, width: 10},
		{name: attrTime, width: 8},
		{name: attrIsDir, width: 0, hidden: true},
	}

	if len(header) != len(wantHeader) {
		t.Fatalf("header has %d attributes, want %d", len(header), len(wantHeader))
	}

	for idx, want := range wantHeader {
		if header[idx].AttrName != want.name {
			t.Errorf("header[%d].AttrName = %q, want %q", idx, header[idx].AttrName, want.name)
		}

		if header[idx].Width != want.width {
			t.Errorf("header[%d].Width = %d, want %d", idx, header[idx].Width, want.width)
		}

		if header[idx].Hidden != want.hidden {
			t.Errorf("header[%d].Hidden = %v, want %v", idx, header[idx].Hidden, want.hidden)
		}

		if header[idx].Flex != want.flex {
			t.Errorf("header[%d].Flex = %v, want %v", idx, header[idx].Flex, want.flex)
		}
	}
}

func assertFileRow(t *testing.T, row app.AttributeList) {
	t.Helper()

	if got := stringAttr(row, attrName); got != "a.txt" {
		t.Errorf("first entry name = %q, want a.txt", got)
	}

	if got := stringAttr(row, attrSize); got != "1" {
		t.Errorf("a.txt size = %q, want 1", got)
	}

	if got := stringAttr(row, attrDate); !datePattern.MatchString(got) {
		t.Errorf("a.txt date = %q, want YYYY-MM-DD", got)
	}

	if got := stringAttr(row, attrTime); !timePattern.MatchString(got) {
		t.Errorf("a.txt time = %q, want HH:MM:SS", got)
	}
}

func assertDirRow(t *testing.T, row app.AttributeList) {
	t.Helper()

	if got := stringAttr(row, attrSize); got != "<DIR>" {
		t.Errorf("sub size = %q, want <DIR>", got)
	}

	value, ok := attrValue(row, attrIsDir)
	if !ok {
		t.Fatalf("sub has no %s attribute", attrIsDir)
	}

	if isDir, isBool := value.(bool); !isBool || !isDir {
		t.Errorf("sub IsDir = %v, want true", value)
	}
}

func TestReadDirSizeValueIsSortable(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "small"), make([]byte, 208), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(dir, "big"), make([]byte, 4000), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	rows, err := filesystem.New().ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	small := sizeOf(t, rows, "small")
	big := sizeOf(t, rows, "big")

	if got := fmt.Sprint(small); got != "208" {
		t.Errorf("small size displays %q, want 208", got)
	}

	if got := fmt.Sprint(big); got != "3.9K" {
		t.Errorf("big size displays %q, want 3.9K", got)
	}

	if small >= big {
		t.Errorf("small size %d should be less than big size %d", small, big)
	}
}

// sizeOf returns the app.Size of the entry named name in a ReadDir result.
func sizeOf(t *testing.T, rows []app.AttributeList, name string) app.Size {
	t.Helper()

	for _, row := range rows[1:] {
		if stringAttr(row, attrName) != name {
			continue
		}

		value, found := attrValue(row, attrSize)
		if !found {
			t.Fatalf("entry %q has no %s attribute", name, attrSize)
		}

		size, isSize := value.(app.Size)
		if !isSize {
			t.Fatalf("entry %q size = %T, want app.Size", name, value)
		}

		return size
	}

	t.Fatalf("entry %q not found", name)

	return 0
}

// stringAttr returns the named attribute's value formatted as a string.
func stringAttr(row app.AttributeList, name string) string {
	value, _ := attrValue(row, name)

	return fmt.Sprintf("%v", value)
}

func TestReadDirEmpty(t *testing.T) {
	t.Parallel()

	rows, err := filesystem.New().ReadDir(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	if rows != nil {
		t.Errorf("expected no rows for empty directory, got %d", len(rows))
	}
}

func TestReadDirSymlinkToDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	err := os.Mkdir(filepath.Join(dir, "real"), 0o700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Symlink("real", filepath.Join(dir, "link"))
	if err != nil {
		t.Skipf("symlinks unsupported: %v", err)
	}

	rows, err := filesystem.New().ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if got := entryIsDir(t, rows, "link"); got != true {
		t.Errorf("symlink to directory IsDir = %v, want true", got)
	}
}

func TestReadDirSymlinkToFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "target.txt"), []byte("x"), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Symlink("target.txt", filepath.Join(dir, "link"))
	if err != nil {
		t.Skipf("symlinks unsupported: %v", err)
	}

	rows, err := filesystem.New().ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if got := entryIsDir(t, rows, "link"); got != false {
		t.Errorf("symlink to file IsDir = %v, want false", got)
	}
}

func TestReadDirBrokenSymlink(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	err := os.Symlink("missing", filepath.Join(dir, "link"))
	if err != nil {
		t.Skipf("symlinks unsupported: %v", err)
	}

	rows, err := filesystem.New().ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if got := entryIsDir(t, rows, "link"); got != false {
		t.Errorf("broken symlink IsDir = %v, want false", got)
	}
}

// entryIsDir returns the IsDir value of the entry named name in a ReadDir
// result, failing the test when the entry is absent or malformed.
func entryIsDir(t *testing.T, rows []app.AttributeList, name string) bool {
	t.Helper()

	for _, row := range rows[1:] {
		rowName, found := attrValue(row, attrName)
		if !found || rowName != name {
			continue
		}

		raw, found := attrValue(row, attrIsDir)
		if !found {
			t.Fatalf("entry %q has no %s attribute", name, attrIsDir)
		}

		value, isBool := raw.(bool)
		if !isBool {
			t.Fatalf("entry %q %s = %T, want bool", name, attrIsDir, raw)
		}

		return value
	}

	t.Fatalf("entry %q not found in listing", name)

	return false
}

func attrValue(row app.AttributeList, name string) (any, bool) {
	for _, attr := range row {
		if attr.AttrName == name {
			return attr.AttrValue, true
		}
	}

	return nil, false
}
