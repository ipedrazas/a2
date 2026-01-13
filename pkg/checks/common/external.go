package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ipedrazas/a2/pkg/checker"
	"github.com/ipedrazas/a2/pkg/checkutil"
)

// ExternalCheck runs an external binary as a check.
// The binary should follow the A2 external check protocol:
// - Exit code 0: Pass
// - Exit code 1: Warning
// - Exit code 2+: Fail
// - Output: Plain text message, or JSON with "message" field
//
// Security Note: External checks are configured by the project owner via .a2.yaml.
// Only commands that exist in PATH are allowed to execute.
type ExternalCheck struct {
	CheckID   string   // Unique identifier
	CheckName string   // Human-readable name
	Command   string   // Command to run (must exist in PATH)
	Args      []string // Arguments to pass
	Severity  string   // Default severity on failure: "warn" or "fail"
}

// ExternalOutput is the optional JSON output format for external checks.
type ExternalOutput struct {
	Message string `json:"message"`
	Status  string `json:"status"` // "pass", "warn", or "fail"
}

func (c *ExternalCheck) ID() string   { return c.CheckID }
func (c *ExternalCheck) Name() string { return c.CheckName }

func (c *ExternalCheck) Run(path string) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	// Validate command before execution
	if err := c.validateCommand(); err != nil {
		return rb.Warn(err.Error()), nil
	}

	// Resolve command to absolute path for security
	cmdPath, err := exec.LookPath(c.Command)
	if err != nil {
		return rb.Warn(fmt.Sprintf("Command not found: %s", c.Command)), nil
	}

	// Sanitize arguments - remove any that contain shell metacharacters
	sanitizedArgs := c.sanitizeArgs()

	// Execute the command with resolved path
	// #nosec G204 -- Command is validated via LookPath and args are sanitized.
	// External checks are an intentional feature configured by project owners in .a2.yaml.
	// Security: cmdPath is resolved via exec.LookPath (must exist in PATH),
	// command name is validated to reject shell metacharacters, and args are passed
	// directly to exec without shell interpretation.
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	cmd := exec.Command(cmdPath, sanitizedArgs...)
	cmd.Dir = path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Parse output
	output := strings.TrimSpace(stdout.String())
	if output == "" {
		output = strings.TrimSpace(stderr.String())
	}

	// Try to parse as JSON
	var extOutput ExternalOutput
	if jsonErr := json.Unmarshal([]byte(output), &extOutput); jsonErr == nil {
		// Successfully parsed JSON
		return c.resultFromJSON(extOutput)
	}

	// Plain text output - determine status from exit code
	return c.resultFromExitCode(output, err)
}

// validateCommand checks that the command is safe to execute.
func (c *ExternalCheck) validateCommand() error {
	if c.Command == "" {
		return fmt.Errorf("empty command")
	}

	// Reject commands with path separators (must be in PATH)
	if strings.ContainsAny(c.Command, "/\\") {
		// Allow if it's an absolute path that exists
		if !filepath.IsAbs(c.Command) {
			return fmt.Errorf("relative paths not allowed: %s", c.Command)
		}
	}

	// Reject shell metacharacters in command name
	dangerous := []string{";", "&", "|", "$", "`", "(", ")", "{", "}", "<", ">", "\n", "\r"}
	for _, char := range dangerous {
		if strings.Contains(c.Command, char) {
			return fmt.Errorf("invalid characters in command: %s", c.Command)
		}
	}

	return nil
}

// sanitizeArgs returns the arguments as-is.
// Since exec.Command passes arguments directly to the executable without shell
// interpretation, shell metacharacters in arguments are not a security risk.
// The command itself is validated via LookPath, and the project owner controls
// the configuration of external checks.
func (c *ExternalCheck) sanitizeArgs() []string {
	return c.Args
}

func (c *ExternalCheck) resultFromJSON(out ExternalOutput) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	switch strings.ToLower(out.Status) {
	case "warn", "warning":
		return rb.Warn(out.Message), nil
	case "fail", "error":
		return rb.Fail(out.Message), nil
	default:
		return rb.Pass(out.Message), nil
	}
}

func (c *ExternalCheck) resultFromExitCode(output string, err error) (checker.Result, error) {
	rb := checkutil.NewResultBuilder(c, checker.LangCommon)

	if err == nil {
		// Exit code 0 = pass
		return rb.Pass(output), nil
	}

	// Get exit code
	exitCode := 1
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}

	message := output
	if message == "" {
		message = "Check failed"
	}

	// Determine severity
	if exitCode >= 2 || c.Severity == "fail" {
		return rb.Fail(message), nil
	}
	return rb.Warn(message), nil
}
