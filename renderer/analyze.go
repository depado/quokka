package renderer

import (
	"os"
	"path/filepath"

	"github.com/Depado/projectmpl/utils"
	"github.com/fatih/color"

	"github.com/Depado/projectmpl/conf"
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

// Analyze is a work in progress function to analyze the template directory
// and gather information about where the configuration files are stored and to
// which templates they should apply.
func Analyze(dir string) {
	var err error
	var root *conf.Root

	if root = GetRootConfig(dir); root == nil {
		utils.FatalPrintln("Couldn't find configuration in template")
	}
	utils.OkPrintln("Root configuration found")
	if err = root.Parse(); err != nil {
		utils.ErrPrintln("Couldn't parse root configuration:", err)
	}
	utils.OkPrintln("Preparing", color.GreenString(root.Name), "-", color.YellowString(root.Version))

	root.ParseVars()

	m := make(map[string]*conf.ConfigFile)
	m[root.File.Dir] = &root.ConfigFile

	// Cycle through to find override configuration files
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == ConfigName && path != root.File.Path {
			cf := conf.NewConfigFile(path, info)
			m[cf.File.Dir] = cf
			utils.OkPrintln("Override Configuration:", color.YellowString(path))
		}
		return nil
	})

	// Cycle through the files
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != ConfigName && info.Name() != ".git" {
			f := conf.NewFile(path, info)
			c := filepath.Dir(path)
			for {
				if v, ok := m[c]; ok {
					v.AddCandidate(f)
					f.Renderers = append(f.Renderers, v)
				}
				if c == root.File.Dir {
					break
				}
				c = filepath.Dir(c)
			}
			conf.AllCandidates = append(conf.AllCandidates, f)
		}
		return nil
	})
}
