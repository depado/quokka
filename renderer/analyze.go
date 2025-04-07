package renderer

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/depado/quokka/conf"
	"github.com/depado/quokka/utils"
)

// ConfigName is the generic name of the file that acts at the configuration
const ConfigName = ".quokka.yml"

// GetRootConfig returns the root configuration that is expected to be at the
// root of the template. Returns nil if the root configuration cannot be found
func GetRootConfig(dir string, ctx conf.InputCtx) *conf.Root {
	exp := filepath.Join(dir, ConfigName)
	info, err := os.Stat(exp)
	if os.IsNotExist(err) {
		return nil
	}
	return conf.NewRootConfig(exp, info, ctx)
}

// HandleRootConfig will find and parse the root configuration. It will then ask
// the user for the variables in the root configuration
func HandleRootConfig(dir string, ctx conf.InputCtx) *conf.Root {
	var err error
	var root *conf.Root

	if root = GetRootConfig(dir, ctx); root == nil {
		utils.FatalPrintln("Couldn't find configuration in template")
	}
	if err = root.Parse(); err != nil {
		utils.FatalPrintln("Couldn't parse root configuration:", err)
	}
	utils.OkPrintln(color.GreenString(root.Name), "-", color.YellowString(root.Version))
	if root.Description != "" {
		utils.OkPrintln(color.CyanString(root.Description))
	}
	root.Prompt()
	return root
}

// Analyze is a work in progress function to analyze the template directory
// and gather information about where the configuration files are stored and to
// which templates they should apply.
func Analyze(dir, output, input string, set []string) {
	var err error
	var ctx conf.InputCtx

	if input != "" {
		if ctx, err = conf.GetInputContext(input); err != nil {
			utils.FatalPrintln("Could not parse input file:", err)
		}
		utils.OkPrintln("Input file", utils.Green.Sprint(input), "found")
	}
	if len(set) > 0 {
		setCtx, err := conf.GetSetContext(set)
		if err != nil {
			utils.FatalPrintln("Could not parse set flags:", err)
		}
		ctx = conf.MergeCtx(ctx, setCtx)
		utils.OkPrintln("Command line set merged in context")
	}

	root := HandleRootConfig(dir, ctx)
	var candidates []*conf.File

	m := make(map[string]*conf.ConfigFile)
	m[root.File.Dir] = &root.ConfigFile

	// Cycle through to find override configuration files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == ConfigName && path != root.File.Path {
			cf := conf.NewConfigFile(path, info, ctx)
			m[cf.File.Dir] = cf
			utils.OkPrintln("Override configuration:", color.YellowString(path))
			if err := cf.Parse(); err != nil {
				utils.FatalPrintln("Couldn't parse configuration:", err)
			}
			cf.Prompt()
		}
		return nil
	})
	if err != nil {
		utils.FatalPrintln("Couldn't read filesystem:", err)
	}

	// Cycle through the files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() != ConfigName && info.Name() != ".git" {
			f := conf.NewFile(path, info, ctx)
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
			candidates = append(candidates, f)
		}
		return nil
	})
	if err != nil {
		utils.FatalPrintln("Couldn't read filesystem:", err)
	}

	for _, f := range candidates {
		if err = f.ParseFrontMatter(); err != nil {
			utils.FatalPrintln("Couldn't parse front matter for file", color.YellowString(f.Path), ":", err)
		}
		if err = f.Render(); err != nil {
			utils.FatalPrintln("Couldn't render template:", err)
		}
	}
}
