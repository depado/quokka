package provider

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type local struct {
	Path string
}

// NewLocalProvider will return a new provider to a local filesystem path
func NewLocalProvider(in string) Provider {
	return local{
		Path: in,
	}
}

func (l local) Fetch() (string, error) {
	return filepath.Join(l.Path, viper.GetString("path")), nil
}

func (local) UsesTmp() bool {
	return false
}

func (local) Name() string {
	return "local"
}
