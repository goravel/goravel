package color

import (
	"bytes"
	"io"
	"os"

	"github.com/pterm/pterm"

	"github.com/goravel/framework/contracts/support"
	"github.com/goravel/framework/support/env"
)

const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
	// FgDefault revert default FG.
	FgDefault Color = 39
)

// Extra foreground color 90 - 97.
const (
	FgDarkGray Color = iota + 90
	FgLightRed
	FgLightGreen
	FgLightYellow
	FgLightBlue
	FgLightMagenta
	FgLightCyan
	FgLightWhite
	// FgGray is an alias of FgDarkGray.
	FgGray Color = 90
)

var (
	info = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.DefaultText,
		Prefix: pterm.Prefix{
			Style: &pterm.Style{pterm.FgBlack, pterm.BgLightWhite},
			Text:  " INFO  ",
		},
		Writer: os.Stdout,
	}
	warning = pterm.Warning
	err     = pterm.Error
	debug   = pterm.Debug
	success = pterm.Success
)

func init() {
	pterm.EnableDebugMessages()
	// Temporarily fix output issue including by https://github.com/pterm/pterm/commit/825931aa7ab264074e6c4045c3bdbca5482c758c
	if env.IsTesting() {
		info.Writer = nil
		warning.Writer = nil
		err.Writer = nil
		debug.Writer = nil
		success.Writer = nil
	}
}

// New Functions to create Printer with specific color
func New(color Color) support.Printer {
	return color
}

func Green() support.Printer {
	return New(FgGreen)
}

func Red() support.Printer {
	return New(FgRed)
}

func Blue() support.Printer {
	return New(FgBlue)
}

func Yellow() support.Printer {
	return New(FgYellow)
}

func Cyan() support.Printer {
	return New(FgCyan)
}

func White() support.Printer {
	return New(FgWhite)
}

func Gray() support.Printer {
	return New(FgGray)
}

func Default() support.Printer {
	return New(FgDefault)
}

func Black() support.Printer {
	return New(FgBlack)
}

func Magenta() support.Printer {
	return New(FgMagenta)
}

type Color uint8

func (c Color) Sprint(a ...any) string {
	return pterm.Color(c).Sprint(a...)
}

func (c Color) Sprintln(a ...any) string {
	return pterm.Color(c).Sprintln(a...)
}

func (c Color) Sprintf(format string, a ...any) string {
	return pterm.Color(c).Sprintf(format, a...)
}

func (c Color) Sprintfln(format string, a ...any) string {
	return pterm.Color(c).Sprintfln(format, a...)
}

func (c Color) Print(a ...any) *support.Printer {
	pterm.Color(c).Print(a...)
	p := support.Printer(c)
	return &p
}

func (c Color) Println(a ...any) *support.Printer {
	pterm.Color(c).Println(a...)
	p := support.Printer(c)
	return &p
}

func (c Color) Printf(format string, a ...any) *support.Printer {
	pterm.Color(c).Printf(format, a...)
	p := support.Printer(c)
	return &p
}

func (c Color) Printfln(format string, a ...any) *support.Printer {
	pterm.Color(c).Printfln(format, a...)
	p := support.Printer(c)
	return &p
}

// Quick use color print message

func Debugf(format string, a ...any) { debug.Printf(format, a...) }

func Debugln(a ...any) { debug.Println(a...) }

func Errorf(format string, a ...any) { err.Printf(format, a...) }

func Errorln(a ...any) { err.Println(a...) }

func Infof(format string, a ...any) { info.Printf(format, a...) }

func Infoln(a ...any) { info.Println(a...) }

func Successf(format string, a ...any) { success.Printf(format, a...) }

func Successln(a ...any) { success.Println(a...) }

// DEPRECATED: Use Warningf instead
func Warnf(format string, a ...any) { warning.Printf(format, a...) }

// DEPRECATED: Use Warningln instead
func Warnln(a ...any) { warning.Println(a...) }

func Warningf(format string, a ...any) { warning.Printf(format, a...) }

func Warningln(a ...any) { warning.Println(a...) }

// CaptureOutput simulates capturing of os.stdout with a buffer and returns what was written to the screen
func CaptureOutput(f func(w io.Writer)) string {
	var outBuf bytes.Buffer
	info.Writer = &outBuf
	warning.Writer = &outBuf
	err.Writer = &outBuf
	debug.Writer = &outBuf
	success.Writer = &outBuf
	pterm.SetDefaultOutput(&outBuf)
	f(&outBuf)

	content := outBuf.String()
	outBuf.Reset()
	return content
}

// Disable disables color output
func Disable() {
	pterm.DisableColor()
}

// Enable enables color output
func Enable() {
	pterm.EnableColor()
}
