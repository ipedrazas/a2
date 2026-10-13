//go:build !windows

package checkutil

import (
	"os/exec"
	"syscall"
)

// setProcessGroup runs cmd in its own process group and makes cancellation
// kill the whole group, not just the direct child.
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		// A negative pid signals the process group.
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}
}
