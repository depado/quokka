package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey"
	"github.com/Depado/projectmpl/utils"
	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v2"
)

// Variable represents a single variable
type Variable struct {
	Type    string   `yaml:"type"`
	Default string   `yaml:"default"`
	Values  []string `yaml:"values"`
	Help    string   `yaml:"help"`
	Result  string
}

// AllCandidates is the full list of candidates
var AllCandidates []*File

// File represents a single file, combining both its path and its os.FileInfo
type File struct {
	Path      string
	Dir       string
	Info      os.FileInfo
	Renderers []*ConfigFile
}

// Config is a configuration that can be applied to a single file (inline conf)
// or to an entire directory
type Config struct {
	Delimiters []string             `yaml:"delimiters"`
	Ignore     bool                 `yaml:"ignore"`
	Variables  map[string]*Variable `yaml:"variables"`
}

// ConfigFile is the combination of File and Config
type ConfigFile struct {
	Config     `yaml:",inline"`
	Candidates []*File
	File       *File
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
	ConfigFile `yaml:",inline"`
	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
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

// ParseVars parses the variables in the root config file
func (r *Root) ParseVars() {
	for name, variable := range r.Variables {
		if len(variable.Values) != 0 {
			survey.SelectQuestionTemplate = surveySelectTemplate
			prompt := &survey.Select{
				Message: fmt.Sprintf("Select a value for %s:", color.YellowString(name)),
				Options: variable.Values,
				Default: variable.Default,
				Help:    variable.Help,
			}
			if err := survey.AskOne(prompt, &variable.Result, nil); err != nil {
				utils.FatalPrintln("Couldn't get an answer:", err)
			}
			utils.OkPrintln("Chose:", color.BlueString(variable.Result), "for", color.YellowString(name), "variable")
		} else {
			prompt := &survey.Input{
				Message: fmt.Sprintf("Value for a %s:", name),
				Default: variable.Default,
				Help:    variable.Help,
			}
			if err := survey.AskOne(prompt, &variable.Result, nil); err != nil {
				utils.FatalPrintln("Couldn't get an answer:", err)
			}
			utils.OkPrintln("Chose:", color.BlueString(variable.Result), "for", color.YellowString(name), "variable")
		}
	}
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
