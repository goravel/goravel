package keyboard

import "syscall"

func closeInput() {
	syscall.CancelIoEx(syscall.Handle(inputTTY.Fd()), nil)
}
