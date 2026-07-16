package conf

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

// Command is a shell command to run in the output directory once the
// template has been rendered
type Command struct {
	Cmd     string `yaml:"cmd"`     // Shell command to run
	Echo    string `yaml:"echo"`    // Message displayed after a successful run
	If      string `yaml:"if"`      // Condition: single variable name or expr-lang expression
	Failure string `yaml:"failure"` // "stop" aborts remaining commands on failure
}

// Run executes the command in the given directory, streaming output to the
// user's terminal
func (c Command) Run(dir string) error {
	ecmd := exec.Command("sh", "-c", c.Cmd) // ponytail: sh -c only, no windows shell support
	ecmd.Dir = dir
	ecmd.Stdin = os.Stdin
	ecmd.Stdout = os.Stdout
	ecmd.Stderr = os.Stderr
	return ecmd.Run()
}

// Include defines an external template to pull in during rendering
type Include struct {
	Source  string `yaml:"source"`  // URL ending in .git or local path
	Path    string `yaml:"path"`    // Inner path within the fetched template (optional)
	Dest    string `yaml:"dest"`    // Sub-directory within the output to render into (default: root)
	If      string `yaml:"if"`      // Condition: single variable name or expr-lang expression
	Confirm *bool  `yaml:"confirm"` // When set, prompt the user with a yes/no question before including
	Prompt  string `yaml:"prompt"`  // Custom message for the confirm prompt (default: "Include <source>?")
}

// Root is a ConfigFile with extra information. It should be located at the root
// of the template
type Root struct {
	ConfigFile  `yaml:",inline"`
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	Includes    []Include `yaml:"includes"`
	After       []Command `yaml:"after"`
}

// Parse will parse the yaml file and store its result in the root config
func (r *Root) Parse() error {
	var err error
	var out []byte

	if out, err = os.ReadFile(r.File.Path); err != nil {
		return err
	}

	return yaml.UnmarshalWithOptions(out, r, yaml.UseOrderedMap())
}

// NewPath adds the path where the file should be rendered according to the root
func (r Root) NewPath(f *File, new string) {
	f.NewPath = filepath.ToSlash(strings.Replace(f.Path, r.File.Dir, new, 1))
}

// NewRootConfig will return the root configuration
func NewRootConfig(path string, file os.FileInfo, ctx InputCtx) *Root {
	return &Root{
		ConfigFile: ConfigFile{
			File: &File{
				Path: path,
				Info: file,
				Dir:  filepath.Dir(path),
			},
			Config: Config{Ctx: ctx},
		},
	}
}
