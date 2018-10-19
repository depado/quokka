package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	yaml "gopkg.in/yaml.v2"
)

// Config is a configuration that can be applied to a single file (inline conf)
// or to an entire directory
type Config struct {
	Delimiters []string             `yaml:"delimiters"`
	Ignore     bool                 `yaml:"ignore"`
	Variables  map[string]*Variable `yaml:"variables"`
	If         string               `yaml:"if"`
}

// PromptVariables will prompt the user for the different variables in the file
func (c *Config) PromptVariables() {
	// Order the variables alphabetically to keep the same order
	var ordered []*Variable
	// ordered := make([]*Variable, len(c.Variables))
	for k, v := range c.Variables {
		if v == nil { // Unconfigured values do have a key but no value
			v = &Variable{Name: k}
		} else {
			v.Name = k
		}
		ordered = append(ordered, v)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Name < ordered[j].Name
	})

	for _, variable := range ordered {
		variable.Prompt()
	}
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
