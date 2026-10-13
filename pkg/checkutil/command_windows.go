//go:build windows

package checkutil

import "os/exec"

// setProcessGroup is a no-op on Windows: exec.CommandContext's default Cancel
// kills the direct child, and WaitDelay bounds any leftover descendants.
func setProcessGroup(_ *exec.Cmd) {}
