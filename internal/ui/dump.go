package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const dumpFileSuffix = ".ansi"

// dumper writes a captured screen frame to a dump file.
type dumper interface {
	Dump(frame string) error
}

// fileDumper writes frames to timestamped files in dir (empty = process
// working directory), appending a numeric suffix when a name is taken.
type fileDumper struct {
	dir string
	now func() time.Time
}

// dumpFilename returns the base name for a frame captured at now, e.g.
// screen-20060102-150405.ansi.
func dumpFilename(now time.Time) string {
	return "screen-" + now.Format("20060102-150405") + dumpFileSuffix
}

func (d *fileDumper) Dump(frame string) error {
	now := d.now
	if now == nil {
		now = time.Now
	}

	path, err := d.nextFreePath(now())
	if err != nil {
		return fmt.Errorf("screen dump: %w", err)
	}

	err = os.WriteFile(path, []byte(frame), 0o600)
	if err != nil {
		return fmt.Errorf("write screen dump %s: %w", path, err)
	}

	return nil
}

// nextFreePath returns the first path in dir whose name is free, trying the
// base timestamped name and then `-1`, `-2`, ... on collision.
func (d *fileDumper) nextFreePath(now time.Time) (string, error) {
	base := dumpFilename(now)

	for i := 0; ; i++ {
		name := base
		if i > 0 {
			name = strings.TrimSuffix(base, dumpFileSuffix) + "-" + strconv.Itoa(i) + dumpFileSuffix
		}

		path := filepath.Join(d.dir, name)

		_, err := os.Stat(path)
		switch {
		case err == nil:
			continue
		case errors.Is(err, os.ErrNotExist):
			return path, nil
		default:
			return "", fmt.Errorf("check screen dump path %s: %w", path, err)
		}
	}
}
