package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/connectors/filesystem"
)

const (
	attrName  = "Name"
	attrIsDir = "IsDir"
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

	wantNames := []string{"Name", "IsDir", "Type"}

	header := rows[0]
	if len(header) != len(wantNames) {
		t.Fatalf("header has %d columns, want %d", len(header), len(wantNames))
	}

	for i, want := range wantNames {
		if header[i].AttrName != want {
			t.Errorf("header[%d].AttrName = %q, want %q", i, header[i].AttrName, want)
		}
	}

	if rows[1][0].AttrValue != "a.txt" {
		t.Errorf("first entry = %v, want a.txt", rows[1][0].AttrValue)
	}

	if rows[2][1].AttrValue != true {
		t.Errorf("sub IsDir = %v, want true", rows[2][1].AttrValue)
	}
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
