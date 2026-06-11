package conf

import (
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

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
}

// Parse will parse the yaml file and store its result in the root config
func (r *Root) Parse() error {
	var err error
	var out []byte

	if out, err = os.ReadFile(r.File.Path); err != nil {
		return err
	}

	return yaml.Unmarshal(out, r)
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
