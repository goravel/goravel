//go:build windows

package process

import (
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/windows"

	"github.com/goravel/framework/errors"
)

// stillActive is a Win32 constant that indicates a process is still running.
// It is not exported by the Go standard library, so we define it here.
const stillActive = 259

// setSysProcAttr configures the process to start in a new process group so we
// can later deliver CTRL_BREAK events to the whole group if needed.
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &windows.SysProcAttr{CreationFlags: windows.CREATE_NEW_PROCESS_GROUP}
}

func running(p *os.Process) bool {
	if p == nil {
		return false
	}

	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(p.Pid))
	if err != nil {
		// If we cannot open the process (access denied or not found), assume not running.
		return false
	}
	defer windows.CloseHandle(h)

	var code uint32
	if err := windows.GetExitCodeProcess(h, &code); err != nil {
		return false
	}
	return code == stillActive
}

func kill(p *os.Process) error {
	if p == nil {
		return errors.ProcessNotStarted
	}
	return p.Kill()
}

func signal(p *os.Process, sig os.Signal) error {
	if p == nil {
		return errors.ProcessNotStarted
	}

	// Map os.Interrupt to CTRL_BREAK_EVENT for the process group when possible
	if sig == os.Interrupt {
		// GenerateConsoleCtrlEvent sends the signal to the process group if
		// the process was created with CREATE_NEW_PROCESS_GROUP.
		// Note: This requires the target to be attached to a console.
		pid := uint32(p.Pid)
		// 0 sends to the process group identified by pid
		if err := windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, pid); err == nil {
			return nil
		}
	}

	return p.Signal(sig)
}

func stop(p *os.Process, done <-chan struct{}, timeout time.Duration, _ ...os.Signal) error {
	if p == nil {
		return errors.ProcessNotStarted
	}

	if !running(p) {
		return nil
	}

	// Try a CTRL_BREAK first for a graceful termination
	_ = signal(p, os.Interrupt)

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		// Force kill if the process didn't exit in time
		return kill(p)
	}
}
