package color

import "github.com/pterm/pterm"

/*
Package color provides utilities for rendering and printing colorized text using color tags.

The color tags are automatically parsed and rendered.

Supported Features:
1. Style Tags:
   Use predefined style tags to format text. Example:
   color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>")

2. Custom Color Attributes:
   Use specific color attributes for more granular control. Example:
   color.Println("<fg=11aa23>he</><bg=120,35,156>llo</>, <fg=167;bg=232>wel</><fg=red>come</>")
*/

// Print renders color tags and prints the provided arguments without a newline.
//
// Example:
//
// color.Print("<suc>Hello</>, <red>World</>")
func Print(a ...any) {
	pterm.Print(a...)
}

// Println renders color tags and prints the provided arguments with a newline.
//
// Example:
//
// color.Println("<cyan>Welcome</>, <green>to</>, <yellow>Go</>")
func Println(a ...any) {
	pterm.Println(a...)
}

// Printf formats a string with the provided arguments, renders color tags, and prints the result.
//
// Example:
//
// color.Printf("<red>Error:</> %s\n", err)
func Printf(format string, a ...any) {
	pterm.Printf(format, a...)
}

// Printfln formats a string with the provided arguments, renders color tags,
// and prints the result with a newline.
//
// Example:
//
// color.Printfln("<success>Success:</> %s", message)
func Printfln(format string, a ...any) {
	pterm.Printfln(format, a...)
}

// Sprint renders color tags and returns the formatted string.
//
// Example:
//
// result := color.Sprint("<blue>Processing...</>")
func Sprint(a ...any) string {
	return pterm.Sprint(a...)
}

// Sprintln renders color tags and returns the formatted string with a newline.
//
// Example:
//
// result := color.Sprintln("<fg=yellow;op=bold>Loading complete.</>")
func Sprintln(a ...any) string {
	return pterm.Sprintln(a...)
}

// Sprintf formats a string with the provided arguments, renders color tags,
// and returns the formatted result.
//
// Example:
//
// result := color.Sprintf("<red>Error:</> %s", err)
func Sprintf(format string, a ...any) string {
	return pterm.Sprintf(format, a...)
}

// Sprintfln formats a string with the provided arguments, renders color tags,
// and returns the formatted result with a newline.
//
// Example:
//
// result := color.Sprintfln("<green>Task completed:</> %s", taskName)
func Sprintfln(format string, a ...any) string {
	return pterm.Sprintfln(format, a...)
}
