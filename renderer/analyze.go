package renderer

import (
	"os"
	"path/filepath"

	"github.com/Depado/projectmpl/conf"
	"github.com/Depado/projectmpl/utils"
	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// ConfigName is the generic name of the file that acts at the configuration
const ConfigName = ".projectmpl.yml"

// GetRootConfig returns the root configuration that is expected to be at the
// root of the template. Returns nil if the root configuration cannot be found
func GetRootConfig(dir string) *conf.Root {
	exp := filepath.Join(dir, ConfigName)
	info, err := os.Stat(exp)
	if os.IsNotExist(err) {
		return nil
	}
	return conf.NewRootConfig(exp, info)
}

// HandleRootConfig will find and parse the root configuration. It will then ask
// the user for the variables in the root configuration
func HandleRootConfig(dir string) *conf.Root {
	var err error
	var root *conf.Root

	if root = GetRootConfig(dir); root == nil {
		utils.FatalPrintln("Couldn't find configuration in template")
	}
	utils.OkPrintln("Root configuration found")
	if err = root.Parse(); err != nil {
		utils.FatalPrintln("Couldn't parse root configuration:", err)
	}
	utils.OkPrintln("Preparing", color.GreenString(root.Name), "-", color.YellowString(root.Version))
	if root.Description != "" {
		utils.OkPrintln(color.BlueString(root.Description))
	}
	root.PromptVariables()
	return root
}

// Analyze is a work in progress function to analyze the template directory
// and gather information about where the configuration files are stored and to
// which templates they should apply.
func Analyze(dir string) {
	var err error
	output := viper.GetString("output")
	root := HandleRootConfig(dir)

	m := make(map[string]*conf.ConfigFile)
	m[root.File.Dir] = &root.ConfigFile

	// Cycle through to find override configuration files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == ConfigName && path != root.File.Path {
			cf := conf.NewConfigFile(path, info)
			m[cf.File.Dir] = cf
			utils.OkPrintln("Override Configuration:", color.YellowString(path))
			if err := cf.Parse(); err != nil {
				utils.FatalPrintln("Couldn't parse configuration:", err)
			}
			cf.PromptVariables()
		}
		return nil
	})
	if err != nil {
		utils.FatalPrintln("Couldn't read filesystem:", err)
	}

	// Cycle through the files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != ConfigName && info.Name() != ".git" {
			f := conf.NewFile(path, info)
			c := filepath.Dir(path)
			for {
				if v, ok := m[c]; ok {
					f.AddRenderer(v)
				}
				if c == root.File.Dir {
					break
				}
				c = filepath.Dir(c)
			}
			root.NewPath(f, output)
			conf.AllCandidates = append(conf.AllCandidates, f)
		}
		return nil
	})
	if err != nil {
		utils.FatalPrintln("Couldn't read filesystem:", err)
	}

	for _, f := range conf.AllCandidates {
		if err = f.ParseFrontMatter(); err != nil {
			utils.FatalPrintln("Couldn't parse front matter for file", color.YellowString(f.Path), ":", err)
		}
		if err = os.MkdirAll(filepath.Dir(f.NewPath), os.ModePerm); err != nil {
			utils.FatalPrintln("Couldn't create directory:", err)
		}
		if err = f.Render(); err != nil {
			utils.FatalPrintln("Couldn't render template:", err)
		}
		utils.OkPrintln("Rendered", color.GreenString(f.NewPath))
	}
}
