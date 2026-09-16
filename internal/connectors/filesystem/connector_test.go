package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/drekunov/gc/internal/connectors/filesystem"
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

	wantNames := []string{"path", "Name", "IsDir", "Type"}

	header := rows[0]
	if len(header) != len(wantNames) {
		t.Fatalf("header has %d columns, want %d", len(header), len(wantNames))
	}

	for i, want := range wantNames {
		if header[i].AttrName != want {
			t.Errorf("header[%d].AttrName = %q, want %q", i, header[i].AttrName, want)
		}
	}

	if rows[1][1].AttrValue != "a.txt" {
		t.Errorf("first entry = %v, want a.txt", rows[1][1].AttrValue)
	}

	if rows[2][2].AttrValue != true {
		t.Errorf("sub IsDir = %v, want true", rows[2][2].AttrValue)
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
