package eris

import (
	"fmt"
	"runtime"
	"strings"
)

// Stack is an array of stack frames stored in a human readable format.
type Stack []StackFrame

// format returns an array of formatted stack frames.
func (s Stack) format(sep string, invert bool) []string {
	var str []string
	for _, f := range s {
		if invert {
			str = append(str, f.format(sep))
		} else {
			str = append([]string{f.format(sep)}, str...)
		}
	}
	return str
}

// StackFrame stores a frame's runtime information in a human readable format.
type StackFrame struct {
	Name string
	File string
	Line int
}

// format returns a formatted stack frame.
func (f *StackFrame) format(sep string) string {
	return fmt.Sprintf("%v%v%v%v%v", f.Name, sep, f.File, sep, f.Line)
}

// caller returns a single stack frame. the argument skip is the number of stack frames
// to ascend, with 0 identifying the caller of Caller.
func caller(skip int) *frame {
	pc, _, _, _ := runtime.Caller(skip)
	var f frame = frame(pc)
	return &f
}

// frame is a single program counter of a stack frame.
type frame uintptr

// pc returns the program counter for a frame.
func (f frame) pc() uintptr {
	return uintptr(f) - 1
}

// get returns a human readable stack frame.
func (f frame) get() StackFrame {
	pc := f.pc()
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()

	i := strings.LastIndex(frame.Function, "/")
	name := frame.Function[i+1:]

	return StackFrame{
		Name: name,
		File: frame.File,
		Line: frame.Line,
	}
}

// callers returns a stack trace. the argument skip is the number of stack frames to skip before recording
// in pc, with 0 identifying the frame for Callers itself and 1 identifying the caller of Callers.
func callers(skip int) *stack {
	const depth = 64
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0 : n-2] // todo: change this to filtering out runtime instead of hardcoding n-2
	return &st
}

// stack is an array of program counters.
type stack []uintptr

// insertPC inserts a wrap error program counter (pc) into the correct place of the root error stack trace.
func (s *stack) insertPC(wrapPCs stack) {
	if len(wrapPCs) == 0 {
		return
	} else if len(wrapPCs) == 1 {
		// append the pc to the end if there's only one
		*s = append(*s, wrapPCs[0])
		return
	}
	for at, f := range *s {
		if f == wrapPCs[0] {
			// break if the stack already contains the pc
			break
		} else if f == wrapPCs[1] {
			// insert the first pc into the stack if the second pc is found
			*s = insert(*s, wrapPCs[0], at)
			break
		}
	}
}

// get returns a human readable stack trace.
func (s *stack) get() []StackFrame {
	var stackFrames []StackFrame

	frames := runtime.CallersFrames(*s)
	for {
		frame, more := frames.Next()
		i := strings.LastIndex(frame.Function, "/")
		name := frame.Function[i+1:]
		stackFrames = append(stackFrames, StackFrame{
			Name: name,
			File: frame.File,
			Line: frame.Line,
		})
		if !more {
			break
		}
	}

	return stackFrames
}

// isGlobal determines if the stack trace represents a global error
func (s *stack) isGlobal() bool {
	frames := s.get()
	for _, f := range frames {
		if strings.ToLower(f.Name) == "runtime.doinit" {
			return true
		}
	}
	return false
}

func insert(s stack, u uintptr, at int) stack {
	// this inserts the pc by breaking the stack into two slices (s[:at] and s[at:])
	return append(s[:at], append([]uintptr{u}, s[at:]...)...)
}
