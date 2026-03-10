//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package keyboard

import (
	"os"

	"github.com/containerd/console"
)

func restoreInput() error {
	if con != nil {
		return con.Reset()
	}
	return nil
}

func initInput() error {
	c, err := console.ConsoleFromFile(stdin)
	if err != nil {
		return err
	}
	con = c

	return nil
}

func openInputTTY() (*os.File, error) {
	f, err := os.Open("/dev/tty")
	if err != nil {
		return nil, err
	}
	return f, nil
}
