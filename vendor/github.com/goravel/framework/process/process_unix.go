//go:build !windows

package process

import (
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"

	"github.com/goravel/framework/errors"
)

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &unix.SysProcAttr{Setpgid: true}
}

func running(p *os.Process) bool {
	if p == nil {
		return false
	}

	// unix.Kill with signal 0 is the common alive check
	err := unix.Kill(p.Pid, 0)
	return err == nil
}

func kill(p *os.Process) error {
	if p == nil {
		return errors.ProcessNotStarted
	}

	// kill the whole process group: negative PID addresses the group.
	// If we can't send to the group, fall back to direct Kill of pid.
	if err := unix.Kill(-p.Pid, unix.SIGKILL); err != nil {
		// fallback: try killing the single process
		return unix.Kill(p.Pid, unix.SIGKILL)
	}

	return nil
}

func signal(p *os.Process, sig os.Signal) error {
	if p == nil {
		return errors.ProcessNotStarted
	}

	s, ok := sig.(unix.Signal)
	if !ok {
		return errors.ProcessUnsupportedSignalType
	}

	pid := p.Pid
	// send to whole process group for consistent behavior
	if err := unix.Kill(-pid, s); err != nil {
		// fallback to single process
		return unix.Kill(pid, s)
	}
	return nil
}

func stop(p *os.Process, done <-chan struct{}, timeout time.Duration, sig ...os.Signal) error {
	if p == nil {
		return errors.ProcessNotStarted
	}

	if !running(p) {
		return nil
	}

	var signalToSend os.Signal = unix.SIGTERM
	if len(sig) > 0 {
		signalToSend = sig[0]
	}

	if err := signal(p, signalToSend); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return nil
		}

		return err
	}

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		if err := signal(p, unix.SIGKILL); err != nil {
			if errors.Is(err, os.ErrProcessDone) {
				return nil
			}

			return err
		}
		return nil
	}
}
