package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner extends the spinner.Spinner struct and adds a few useful methods
type Spinner struct {
	*spinner.Spinner
}

// NewSpinner configures and starts a new spinner
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" %s…", message)
	s.Color("green") // nolint: errcheck
	s.Start()
	return &Spinner{s}
}

// ErrStop changes the final message of a Spinner (error) and stops it right
// away
func (is *Spinner) ErrStop(message string, opts ...interface{}) {
	is.FinalMSG = fmt.Sprintln(append([]interface{}{ErrPrefix, message}, opts...)...)
	is.Stop()
}

// DoneStop changes the final message of a Spinner (ok) and stops it right away
func (is *Spinner) DoneStop(message string, opts ...interface{}) {
	is.FinalMSG = fmt.Sprintln(append([]interface{}{OkPrefix, message}, opts...)...)
	is.Stop()
}
