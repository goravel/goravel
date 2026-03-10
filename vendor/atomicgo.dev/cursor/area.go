package cursor

import (
	"os"
	"strings"
)

// Area displays content which can be updated on the fly.
// You can use this to create live output, charts, dropdowns, etc.
type Area struct {
	height     int
	writer     Writer
	cursor     *Cursor
	cursorPosY int
}

// NewArea returns a new Area.
func NewArea() Area {
	return Area{
		height:     0,
		writer:     os.Stdout,
		cursor:     cursor,
		cursorPosY: 0,
	}
}

// WithWriter sets the custom writer.
func (area Area) WithWriter(writer Writer) Area {
	area.writer = writer
	area.cursor = area.cursor.WithWriter(writer)

	return area
}

// Clear clears the content of the Area.
func (area *Area) Clear() {
	// Initialize writer if not done yet
	if area.writer == nil {
		area.writer = os.Stdout
	}

	if area.height > 0 {
		area.Bottom()
		area.ClearLinesUp(area.height)
		area.StartOfLine()
	} else {
		area.StartOfLine()
		area.cursor.ClearLine()
	}
}

// Update overwrites the content of the Area and adjusts its height based on content.
func (area *Area) Update(content string) {
	area.Clear()
	area.writeArea(content)
	area.cursorPosY = 0
	area.height = strings.Count(content, "\n")
}

// Up moves the cursor of the area up one line.
func (area *Area) Up(n int) {
	if n > 0 {
		if area.cursorPosY+n > area.height {
			n = area.height - area.cursorPosY
		}

		area.cursor.Up(n)
		area.cursorPosY += n
	}
}

// Down moves the cursor of the area down one line.
func (area *Area) Down(n int) {
	if n > 0 {
		if area.cursorPosY-n < 0 {
			n = area.height - area.cursorPosY
		}

		area.cursor.Down(n)
		area.cursorPosY -= n
	}
}

// Bottom moves the cursor to the bottom of the terminal.
// This is done by calculating how many lines were moved by Up and Down.
func (area *Area) Bottom() {
	if area.cursorPosY > 0 {
		area.Down(area.cursorPosY)
		area.cursorPosY = 0
	}
}

// Top moves the cursor to the top of the area.
// This is done by calculating how many lines were moved by Up and Down.
func (area *Area) Top() {
	if area.cursorPosY < area.height {
		area.Up(area.height - area.cursorPosY)
		area.cursorPosY = area.height
	}
}

// StartOfLine moves the cursor to the start of the current line.
func (area *Area) StartOfLine() {
	area.cursor.HorizontalAbsolute(0)
}

// StartOfLineDown moves the cursor down by n lines, then moves to cursor to the start of the line.
func (area *Area) StartOfLineDown(n int) {
	area.Down(n)
	area.StartOfLine()
}

// StartOfLineUp moves the cursor up by n lines, then moves to cursor to the start of the line.
func (area *Area) StartOfLineUp(n int) {
	area.Up(n)
	area.StartOfLine()
}

// UpAndClear moves the cursor up by n lines, then clears the line.
func (area *Area) UpAndClear(n int) {
	area.Up(n)
	area.cursor.ClearLine()
}

// DownAndClear moves the cursor down by n lines, then clears the line.
func (area *Area) DownAndClear(n int) {
	area.Down(n)
	area.cursor.ClearLine()
}

// Move moves the cursor relative by x and y.
func (area *Area) Move(x, y int) {
	if x > 0 {
		area.cursor.Right(x)
	} else if x < 0 {
		area.cursor.Left(-x)
	}

	if y > 0 {
		area.Up(y)
	} else if y < 0 {
		area.Down(-y)
	}
}

// ClearLinesUp clears n lines upwards from the current position and moves the cursor.
func (area *Area) ClearLinesUp(n int) {
	area.StartOfLine()
	area.cursor.ClearLine()

	for i := 0; i < n; i++ {
		area.UpAndClear(1)
	}
}

// ClearLinesDown clears n lines downwards from the current position and moves the cursor.
func (area *Area) ClearLinesDown(n int) {
	area.StartOfLine()
	area.cursor.ClearLine()

	for i := 0; i < n; i++ {
		area.DownAndClear(1)
	}
}
