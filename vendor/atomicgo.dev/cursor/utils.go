package cursor

import (
	"io"
	"os"
)

//
// Helpers for global cursor handling on os.Stdout
//

var autoheight int
var cursor = &Cursor{writer: os.Stdout}

// Writer is an expanded io.Writer interface with a file descriptor.
type Writer interface {
	io.Writer
	Fd() uintptr
}

// SetTarget sets to output target of the default curser to the
// provided cursor.Writer (wrapping io.Writer).
func SetTarget(w Writer) {
	cursor = cursor.WithWriter(w)
}

// Up moves the cursor n lines up relative to the current position.
func Up(n int) {
	cursor.Up(n)
	autoheight += n
}

// Down moves the cursor n lines down relative to the current position.
func Down(n int) {
	cursor.Down(n)

	if autoheight > 0 {
		autoheight -= n
	}
}

// Right moves the cursor n characters to the right relative to the current position.
func Right(n int) {
	cursor.Right(n)
}

// Left moves the cursor n characters to the left relative to the current position.
func Left(n int) {
	cursor.Left(n)
}

// HorizontalAbsolute moves the cursor to n horizontally.
// The position n is absolute to the start of the line.
func HorizontalAbsolute(n int) {
	cursor.HorizontalAbsolute(n)
}

// Show the cursor if it was hidden previously.
// Don't forget to show the cursor at least at the end of your application.
// Otherwise the user might have a terminal with a permanently hidden cursor, until they reopen the terminal.
func Show() {
	cursor.Show()
}

// Hide the cursor.
// Don't forget to show the cursor at least at the end of your application with Show.
// Otherwise the user might have a terminal with a permanently hidden cursor, until they reopen the terminal.
func Hide() {
	cursor.Hide()
}

// ClearLine clears the current line and moves the cursor to it's start position.
func ClearLine() {
	cursor.ClearLine()
}

// Clear clears the current position and moves the cursor to the left.
func Clear() {
	cursor.Clear()
}

// Bottom moves the cursor to the bottom of the terminal.
// This is done by calculating how many lines were moved by Up and Down.
func Bottom() {
	if autoheight > 0 {
		Down(autoheight)
		StartOfLine()

		autoheight = 0
	}
}

// StartOfLine moves the cursor to the start of the current line.
func StartOfLine() {
	HorizontalAbsolute(0)
}

// StartOfLineDown moves the cursor down by n lines, then moves to cursor to the start of the line.
func StartOfLineDown(n int) {
	Down(n)
	StartOfLine()
}

// StartOfLineUp moves the cursor up by n lines, then moves to cursor to the start of the line.
func StartOfLineUp(n int) {
	Up(n)
	StartOfLine()
}

// UpAndClear moves the cursor up by n lines, then clears the line.
func UpAndClear(n int) {
	Up(n)
	ClearLine()
}

// DownAndClear moves the cursor down by n lines, then clears the line.
func DownAndClear(n int) {
	Down(n)
	ClearLine()
}

// Move moves the cursor relative by x and y.
func Move(x, y int) {
	if x > 0 {
		Right(x)
	} else if x < 0 {
		x *= -1
		Left(x)
	}

	if y > 0 {
		Up(y)
	} else if y < 0 {
		y *= -1
		Down(y)
	}
}

// ClearLinesUp clears n lines upwards from the current position and moves the cursor.
func ClearLinesUp(n int) {
	for i := 0; i < n; i++ {
		UpAndClear(1)
	}
}

// ClearLinesDown clears n lines downwards from the current position and moves the cursor.
func ClearLinesDown(n int) {
	for i := 0; i < n; i++ {
		DownAndClear(1)
	}
}
