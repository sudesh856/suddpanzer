package reporter

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func Print(requests int64, rps float64, p99 int64, errors int64) {
	line := fmt.Sprintf("Requests: %d | RPS: %.0f | p99: %dms | Errors: %d",
		requests, rps, p99, errors)

	// get terminal width
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 80 // fallback to 80
	}

	// truncate if line is too long for terminal
	if len(line) > width {
		line = line[:width-3] + "..."
	}

	fmt.Printf("\r%-*s", width, line)
}