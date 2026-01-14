package output

import (
	"fmt"
	"os"
	"sync"
)

// ProgressReporter handles progress display to stderr.
type ProgressReporter struct {
	mu          sync.Mutex
	lastMessage string
}

// NewProgressReporter creates a new progress reporter.
func NewProgressReporter() *ProgressReporter {
	return &ProgressReporter{}
}

// Update displays the current progress, overwriting the previous line.
func (p *ProgressReporter) Update(completed, total int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Clear previous line by writing carriage return and spaces
	if p.lastMessage != "" {
		fmt.Fprintf(os.Stderr, "\r%s\r", clearLine(len(p.lastMessage)))
	}

	// Write new progress
	msg := fmt.Sprintf("Running checks... (%d/%d completed)", completed, total)
	fmt.Fprint(os.Stderr, msg)
	p.lastMessage = msg
}

// Done clears the progress line to prepare for final output.
func (p *ProgressReporter) Done() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.lastMessage != "" {
		// Clear the line completely
		fmt.Fprintf(os.Stderr, "\r%s\r", clearLine(len(p.lastMessage)))
		p.lastMessage = ""
	}
}

// clearLine returns a string of spaces to clear N characters.
func clearLine(n int) string {
	spaces := make([]byte, n)
	for i := range spaces {
		spaces[i] = ' '
	}
	return string(spaces)
}
