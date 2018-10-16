package renderer

import (
	"os"
	"path/filepath"

	"github.com/Depado/projectmpl/utils"
	"github.com/fatih/color"
)

// Analyze is a work in progress function to analyze the template directory
// and gather information about where the configuration files are stored and to
// which templates they should apply. TODO: Think about how to create a config
// hierarchy
func Analyze(dir string) {
	var rootcfg string
	m := make(map[string]string)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Ignore directories for the first pass
		if info.IsDir() {
			return nil
		}
		if info.Name() == ".projectmpl.yml" {
			m[filepath.Dir(path)] = path
			if rootcfg == "" {
				utils.OkPrintln("Root Configuration:", color.YellowString(path))
				rootcfg = path
				return nil
			}
			utils.OkPrintln("Override Configuration:", color.YellowString(path))
			return nil
		}
		return nil
	})
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.Name() != ".projectmpl.yml" {
			c := filepath.Dir(path)
			if v, ok := m[c]; ok {
				if v != rootcfg {
					utils.OkPrintln("Candidate:", color.GreenString(path), "rendered by", color.RedString(rootcfg), "and", color.RedString(v))
				} else {
					utils.OkPrintln("Candidate:", color.GreenString(path), "rendered by", color.RedString(v))
				}
				return nil
			}
			c = filepath.Dir(c)
			for {
				if v, ok := m[c]; ok {
					utils.OkPrintln("Candidate:", color.GreenString(path), "rendered by", color.RedString(v))
					return nil
				}
				c = filepath.Dir(path)
			}
		}
		return nil
	})
	utils.OkPrintln("Root Config:", color.GreenString(rootcfg))
}
