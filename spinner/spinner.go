// Package spinner provides a terminal spinner/loader
package spinner

import (
	"fmt"
	"time"

	"github.com/bharath-srinivas/aws-go/colors"
)

// constants for spinner color wrapping
const (
	escape = "\x1b"
	reset  = 0
)

// Spinner config.
type Spinner struct {
	Delay	time.Duration	// Delay is the speed of the spinner
	Prefix	string 			// Prefix string that will be prepended to the spinner
	Color	colors.Color	// Color of the spinner
	active	bool			// status of the spinner
	spinner	[]string		// character set to be used for the spinner
	stop	chan bool		// channel to send stop signal to the spinner
}

// Default is a wrapper around NewSpinner with predefined spinner color and time delay options.
func Default(p string) *Spinner {
	return NewSpinner(p, colors.Cyan, 100 * time.Millisecond)
}

// NewSpinner returns a pointer to Spinner instance with provided options.
func NewSpinner(prefix string, color colors.Color, d time.Duration) *Spinner {
	return &Spinner{
		Delay:		d,
		Prefix:		prefix,
		Color:		color,
		active:		false,
		spinner:	[]string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"},
		stop:		make(chan bool, 1),
	}
}

// Start activates the spinner.
func (s *Spinner) Start() {
	if s.active {
		return
	}

	s.active = true

	go func() {
		for {
			for i := 0; i < len(s.spinner); i++ {
				select {
				case <-s.stop:
					return
				default:
					fmt.Printf("%s%s \r", s.Prefix, s.color(s.spinner[i]))
					time.Sleep(s.Delay)
				}
			}
		}
	}()
}

// Stop deactivates the spinner.
func (s *Spinner) Stop() {
	if s.active {
		s.active = false
		s.stop <- true
	}
}

// color wraps the given string with given color and returns the colored string
func (s *Spinner) color(str string) string {
	if s.Color == 0 {
		return str
	}
	prefix := fmt.Sprintf("%s[%dm", escape, s.Color)
	suffix := fmt.Sprintf("%s[%dm", escape, reset)
	return prefix + str + suffix
}