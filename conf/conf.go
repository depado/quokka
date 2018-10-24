package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Config is a configuration that can be applied to a single file (inline conf)
// or to an entire directory
type Config struct {
	Delimiters []string   `yaml:"delimiters"`
	Copy       *bool      `yaml:"copy"`
	Ignore     *bool      `yaml:"ignore"`
	Variables  *Variables `yaml:"variables"`
	If         string     `yaml:"if"`
}

// ConfigFile is the combination of File and Config
type ConfigFile struct {
	Config `yaml:",inline"`
	File   *File `yaml:"-"`
}

// Parse parses the configuration file
func (c *ConfigFile) Parse() error {
	var err error
	var out []byte

	if out, err = ioutil.ReadFile(c.File.Path); err != nil {
		return err
	}

	return yaml.Unmarshal(out, c)
}

// NewConfigFile returns a new configfile
func NewConfigFile(path string, file os.FileInfo) *ConfigFile {
	return &ConfigFile{File: &File{Path: path, Info: file, Dir: filepath.Dir(path)}}
}

// NewFile returns a new file
func NewFile(path string, file os.FileInfo) *File {
	return &File{Path: path, Info: file, Dir: filepath.Dir(path)}
}
