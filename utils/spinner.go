package utils

import (
	"fmt"
	"strings"

	"github.com/depado/gorich"
	"github.com/depado/gorich/live"
)

type Spinner struct {
	s *live.ActiveSpinner
}

func NewSpinner(message string) *Spinner {
	return &Spinner{s: live.StartSpinner(message)}
}

func (s *Spinner) ErrStop(message string, opts ...any) {
	s.s.Stop()
	gorich.Println("[red]✗[/] " + join(message, opts...))
}

func (s *Spinner) DoneStop(message string, opts ...any) {
	s.s.Stop()
	gorich.Println("[green]✓[/] " + join(message, opts...))
}

func join(message string, opts ...any) string {
	parts := make([]string, 0, len(opts)+1)
	parts = append(parts, message)
	for _, o := range opts {
		parts = append(parts, fmt.Sprint(o))
	}
	return strings.Join(parts, " ")
}
