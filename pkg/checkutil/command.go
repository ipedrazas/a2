package checkutil

import (
	"context"
	"os/exec"
	"time"
)

// waitDelay bounds how long Wait blocks for a cancelled command's output pipes
// to close, so a grandchild that outlives the kill cannot hang the check.
const waitDelay = 5 * time.Second

// Command builds an *exec.Cmd bound to ctx. When ctx is cancelled (timeout or
// fail-fast) the command and, on Unix, every process it spawned are killed, so
// a timed-out `go test` or `npx jest` does not keep running in the background.
// Use this instead of exec.Command inside checks.
func Command(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...) // #nosec G204 -- command execution is by design
	setProcessGroup(cmd)
	cmd.WaitDelay = waitDelay
	return cmd
}
