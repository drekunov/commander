package app

import (
	"fmt"
	"strconv"
)

// Size is a file size in bytes. Its string form is human-readable, while the
// underlying value compares numerically.
type Size int64

// String renders the size with a binary unit and one decimal place, for example
// "512", "4.0K", or "1.2M".
func (s Size) String() string {
	const unit = 1024

	if s < unit {
		return strconv.FormatInt(int64(s), 10)
	}

	div, exp := int64(unit), 0
	for n := int64(s) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f%c", float64(s)/float64(div), "KMGTPE"[exp])
}

type Attribute struct {
	AttrName  string
	AttrValue any

	// Hidden marks an attribute that supplies data to the panel without being
	// rendered as a table column.
	Hidden bool

	// Width is the column width in cells declared on a header attribute; a
	// non-positive value means the panel's default width.
	Width int

	// Flex marks a header attribute whose column fills the width left after the
	// fixed columns and padding, instead of using its declared width.
	Flex bool
}

type AttributeList []Attribute

type Connector interface {
	Name() string
	ReadDir(path string) ([]AttributeList, error)
	ReadFile(path string) ([]byte, error)
}
