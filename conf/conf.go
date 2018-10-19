package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// Config is a configuration that can be applied to a single file (inline conf)
// or to an entire directory
type Config struct {
	Delimiters []string             `yaml:"delimiters"`
	Ignore     bool                 `yaml:"ignore"`
	Variables  map[string]*Variable `yaml:"variables"`
}

// PromptVariables will prompt the user for the different variables in the file
func (c *Config) PromptVariables() {
	// Order the variables alphabetically to keep the same order
	var ordered []*Variable
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
	Config     `yaml:",inline"`
	Candidates []*File `yaml:"-"`
	File       *File   `yaml:"-"`
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

// AddCandidate will add a candidate
func (c *ConfigFile) AddCandidate(f *File) {
	c.Candidates = append(c.Candidates, f)
}

// AddCandidateFromPath will add a candidate from path and file info
func (c *ConfigFile) AddCandidateFromPath(path string, info os.FileInfo) {
	c.Candidates = append(c.Candidates, NewFile(path, info))
}

// Root is a ConfigFile with extra information. It should be located at the root
// of the template
type Root struct {
	ConfigFile  `yaml:",inline"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// Parse will parse the yaml file and store its result in the root config
func (r *Root) Parse() error {
	var err error
	var out []byte

	if out, err = ioutil.ReadFile(r.ConfigFile.File.Path); err != nil {
		return err
	}

	return yaml.Unmarshal(out, r)
}

// NewPath adds the path where the file should be rendered according to the root
func (r Root) NewPath(f *File, new string) {
	f.NewPath = strings.Replace(f.Path, r.File.Dir, new, 1)
}

// NewRootConfig will return the root configuration
func NewRootConfig(path string, file os.FileInfo) *Root {
	return &Root{ConfigFile: ConfigFile{File: &File{Path: path, Info: file, Dir: filepath.Dir(path)}}}
}

// NewConfigFile returns a new configfile
func NewConfigFile(path string, file os.FileInfo) *ConfigFile {
	return &ConfigFile{File: &File{Path: path, Info: file, Dir: filepath.Dir(path)}}
}

// NewFile returns a new file
func NewFile(path string, file os.FileInfo) *File {
	return &File{Path: path, Info: file, Dir: filepath.Dir(path)}
}
