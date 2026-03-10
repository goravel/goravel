//go:build windows
// +build windows

package cursor

import (
	"syscall"
	"unsafe"
)

// Up moves the cursor n lines up relative to the current position.
func (c *Cursor) Up(n int) {
	c.move(0, -n)
}

// Down moves the cursor n lines down relative to the current position.
func (c *Cursor) Down(n int) {
	c.move(0, n)
}

// Right moves the cursor n characters to the right relative to the current position.
func (c *Cursor) Right(n int) {
	c.move(n, 0)
}

// Left moves the cursor n characters to the left relative to the current position.
func (c *Cursor) Left(n int) {
	c.move(-n, 0)
}

func (c *Cursor) move(x int, y int) {
	handle := syscall.Handle(c.writer.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = csbi.cursorPosition.x + short(x)
	cursor.y = csbi.cursorPosition.y + short(y)

	_, _, _ = procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

// HorizontalAbsolute moves the cursor to n horizontally.
// The position n is absolute to the start of the line.
func (c *Cursor) HorizontalAbsolute(n int) {
	handle := syscall.Handle(c.writer.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var cursor coord
	cursor.x = short(n)
	cursor.y = csbi.cursorPosition.y

	if csbi.size.x < cursor.x {
		cursor.x = csbi.size.x
	}

	_, _, _ = procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}

// Show the cursor if it was hidden previously.
// Don't forget to show the cursor at least at the end of your application.
// Otherwise the user might have a terminal with a permanently hidden cursor, until he reopens the terminal.
func (c *Cursor) Show() {
	handle := syscall.Handle(c.writer.Fd())

	var cci consoleCursorInfo
	_, _, _ = procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	cci.visible = 1

	_, _, _ = procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

// Hide the cursor.
// Don't forget to show the cursor at least at the end of your application with Show.
// Otherwise the user might have a terminal with a permanently hidden cursor, until he reopens the terminal.
func (c *Cursor) Hide() {
	handle := syscall.Handle(c.writer.Fd())

	var cci consoleCursorInfo
	_, _, _ = procGetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
	cci.visible = 0

	_, _, _ = procSetConsoleCursorInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&cci)))
}

// ClearLine clears the current line and moves the cursor to its start position.
func (c *Cursor) ClearLine() {
	handle := syscall.Handle(c.writer.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var w uint32
	var x short
	cursor := csbi.cursorPosition
	x = csbi.size.x
	_, _, _ = procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(x), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))
}

// Clear clears the current position and moves the cursor to the left.
func (c *Cursor) Clear() {
	handle := syscall.Handle(c.writer.Fd())

	var csbi consoleScreenBufferInfo
	_, _, _ = procGetConsoleScreenBufferInfo.Call(uintptr(handle), uintptr(unsafe.Pointer(&csbi)))

	var w uint32
	cursor := csbi.cursorPosition
	_, _, _ = procFillConsoleOutputCharacter.Call(uintptr(handle), uintptr(' '), uintptr(1), uintptr(*(*int32)(unsafe.Pointer(&cursor))), uintptr(unsafe.Pointer(&w)))

	if cursor.x > 0 {
		cursor.x--
	}
	_, _, _ = procSetConsoleCursorPosition.Call(uintptr(handle), uintptr(*(*int32)(unsafe.Pointer(&cursor))))
}
