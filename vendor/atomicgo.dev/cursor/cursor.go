package cursor

import (
	"os"
)

// Cursor displays content which can be updated on the fly.
// You can use this to create live output, charts, dropdowns, etc.
type Cursor struct {
	writer Writer
}

// NewCursor creates a new Cursor instance writing to os.Stdout.
func NewCursor() *Cursor {
	return &Cursor{writer: os.Stdout}
}

// WithWriter allows for any arbitrary Writer to be used
// for cursor movement abstracted.
func (c *Cursor) WithWriter(w Writer) *Cursor {
	if w != nil {
		c.writer = w
	}

	return c
}
