//go:build !windows
// +build !windows

package cursor

import (
	"fmt"
)

// Up moves the cursor n lines up relative to the current position.
func (c *Cursor) Up(n int) {
	if n > 0 {
		fmt.Fprintf(c.writer, "\x1b[%dA", n)
	}
}

// Down moves the cursor n lines down relative to the current position.
func (c *Cursor) Down(n int) {
	if n > 0 {
		fmt.Fprintf(c.writer, "\x1b[%dB", n)
	}
}

// Right moves the cursor n characters to the right relative to the current position.
func (c *Cursor) Right(n int) {
	if n > 0 {
		fmt.Fprintf(c.writer, "\x1b[%dC", n)
	}
}

// Left moves the cursor n characters to the left relative to the current position.
func (c *Cursor) Left(n int) {
	if n > 0 {
		fmt.Fprintf(c.writer, "\x1b[%dD", n)
	}
}

// HorizontalAbsolute moves the cursor to n horizontally.
// The position n is absolute to the start of the line.
func (c *Cursor) HorizontalAbsolute(n int) {
	n++ // Moves the line to the character after n
	fmt.Fprintf(c.writer, "\x1b[%dG", n)
}

// Show the cursor if it was hidden previously.
// Don't forget to show the cursor at least at the end of your application.
// Otherwise the user might have a terminal with a permanently hidden cursor, until they reopen the terminal.
func (c *Cursor) Show() {
	fmt.Fprint(c.writer, "\x1b[?25h")
}

// Hide the cursor.
// Don't forget to show the cursor at least at the end of your application with Show.
// Otherwise the user might have a terminal with a permanently hidden cursor, until they reopen the terminal.
func (c *Cursor) Hide() {
	fmt.Fprintf(c.writer, "\x1b[?25l")
}

// ClearLine clears the current line and moves the cursor to it's start position.
func (c *Cursor) ClearLine() {
	fmt.Fprintf(c.writer, "\x1b[2K")
}

// Clear clears the current position and moves the cursor to the left.
func (c *Cursor) Clear() {
	fmt.Fprintf(c.writer, "\x1b[K")
}
