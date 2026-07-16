package conf

import (
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// Config is a configuration that can be applied to a single file (inline conf)
// or to an entire directory
type Config struct {
	Delimiters []string   `yaml:"delimiters"`
	Copy       *bool      `yaml:"copy"`
	Ignore     *bool      `yaml:"ignore"`
	Ignores    []string   `yaml:"ignores"`
	Variables  *Variables `yaml:"variables"`
	If         string     `yaml:"if"`
	Ctx        InputCtx   `yaml:"-"`
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

	if out, err = os.ReadFile(c.File.Path); err != nil {
		return err
	}

	return yaml.UnmarshalWithOptions(out, c, yaml.UseOrderedMap())
}

// Prompt will prompt the necessary prompts without displaying the ones for
// variables that were already defined in the input file
func (c *ConfigFile) Prompt(builtins map[string]any) {
	if c.Variables != nil {
		c.Variables.FillPrompt("", c.Ctx, builtins)
	}
}

// NewConfigFile returns a new configfile
func NewConfigFile(path string, file os.FileInfo, ctx InputCtx) *ConfigFile {
	return &ConfigFile{File: &File{Path: path, Info: file, Dir: filepath.Dir(path)}, Config: Config{Ctx: ctx}}
}

// NewFile returns a new file
func NewFile(path string, file os.FileInfo, ctx InputCtx) *File {
	return &File{Path: path, Info: file, Dir: filepath.Dir(path), Ctx: ctx}
}
