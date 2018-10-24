package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Depado/projectmpl/utils"
	yaml "gopkg.in/yaml.v2"
)

// Root is a ConfigFile with extra information. It should be located at the root
// of the template
type Root struct {
	ConfigFile  `yaml:",inline"`
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	After       []Command `yaml:"after"`
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

// ExecuteCommands will execute the commands in the newly rendered directory
func (r Root) ExecuteCommands(dir string) {
	var err error

	if err = os.Chdir(dir); err != nil {
		utils.FatalPrintln("Couldn't change directory:", err)
	}

	for _, cmd := range r.After {
		if cmd.Cmd != "" {
			if cmd.If != "" && r.Variables != nil {
				if v, ok := r.Variables.m[cmd.If]; ok && v.True() {
					cmd.Run()
				}
			} else {
				cmd.Run()
			}
		}
	}
}

// NewPath adds the path where the file should be rendered according to the root
func (r Root) NewPath(f *File, new string) {
	f.NewPath = strings.Replace(f.Path, r.File.Dir, new, 1)
}

// NewRootConfig will return the root configuration
func NewRootConfig(path string, file os.FileInfo) *Root {
	return &Root{ConfigFile: ConfigFile{File: &File{Path: path, Info: file, Dir: filepath.Dir(path)}}}
}
