package support

type Printer interface {
	Sprint(a ...any) string
	Sprintln(a ...any) string
	Sprintf(format string, a ...any) string
	Sprintfln(format string, a ...any) string
	Print(a ...any) *Printer
	Println(a ...any) *Printer
	Printf(format string, a ...any) *Printer
	Printfln(format string, a ...any) *Printer
}
