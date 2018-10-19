package conf

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Depado/projectmpl/utils"
	yaml "gopkg.in/yaml.v2"
)

// Command is a simple command to be executed
type Command struct {
	Cmd     string `yaml:"cmd"`
	Failure string `yaml:"failure"`
	Output  bool   `yaml:"output"`
	Echo    string `yaml:"echo"`
	If      string `yaml:"if"`
}

// Run will run a single command
func (c Command) Run() {
	var err error
	var output []byte
	var main string
	var args []string

	parts := strings.Fields(c.Cmd)
	if len(parts) == 0 {
		return
	} else if len(parts) == 1 {
		main = parts[0]
	} else if len(parts) > 1 {
		main = parts[0]
		args = parts[1:]
	}

	var ecmd *exec.Cmd
	if len(args) > 0 {
		ecmd = exec.Command(main, args...)
	} else {
		ecmd = exec.Command(main)
	}
	if output, err = ecmd.Output(); err != nil {
		if c.Failure == "stop" {
			utils.FatalPrintln("Couldn't run after command:", err)
		}
		utils.ErrPrintln("Couldn't execute command, ignoring:", err)
		return
	}
	if c.Output {
		out := strings.Split(string(output), "\n")
		for i := 0; i < len(out)-1; i++ {
			utils.OkPrintln(out[i])
		}
	}
	if c.Echo != "" {
		utils.OkPrintln(c.Echo)
	}
}

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
				if v, ok := r.Variables[cmd.If]; ok {
					if v.Confirm != nil && *v.Confirm || v.Result != "" {
						cmd.Run()
					}
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
