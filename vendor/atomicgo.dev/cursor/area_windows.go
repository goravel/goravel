//go:build windows
// +build windows

package cursor

import (
	"fmt"
)

// writeArea is a helper for platform dependant output.
// For Windows newlines '\n' in the content are replaced by '\r\n'
func (area *Area) writeArea(content string) {
	last := ' '
	for _, r := range content {
		if r == '\n' && last != '\r' {
			fmt.Fprint(area.writer, "\r\n")
			continue
		}
		fmt.Fprint(area.writer, string(r))
	}
}
