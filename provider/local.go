package provider

import "path/filepath"

type local struct {
	Path      string
	InnerPath string
}

// NewLocalProvider will return a new provider to a local filesystem path
func NewLocalProvider(in, path string) Provider {
	return local{
		Path:      in,
		InnerPath: path,
	}
}

func (l local) Fetch() (string, error) {
	return filepath.Join(l.Path, l.InnerPath), nil
}

func (local) UsesTmp() bool {
	return false
}

func (local) Name() string {
	return "local"
}
