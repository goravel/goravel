//go:build !windows
// +build !windows

package cursor

import (
	"fmt"
)

// Update overwrites the content of the Area and adjusts its height based on content.
func (area *Area) writeArea(content string) {
	fmt.Fprint(area.writer, content)
}
