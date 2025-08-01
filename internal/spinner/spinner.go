package spinner

import (
	"fmt"
	"time"
)

// Spinner represents a loading spinner
type Spinner struct {
	done chan bool
}

// New creates a new spinner instance
func New() *Spinner {
	return &Spinner{
		done: make(chan bool),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start(model string) {
	go s.run(model)
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.done <- true
}

func (s *Spinner) run(model string) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0

	for {
		select {
		case <-s.done:
			// Clear the entire line and move to next line
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r🤔 %s Thinking with %s model...", spinner[i], model)
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(spinner)
		}
	}
}
