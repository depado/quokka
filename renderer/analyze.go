package renderer

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/depado/quokka/conf"
	"github.com/depado/quokka/provider"
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
func HandleRootConfig(dir string, ctx conf.InputCtx, builtins map[string]any) *conf.Root {
	var err error
	var root *conf.Root

	if root = GetRootConfig(dir, ctx); root == nil {
		utils.FatalPrintln("Couldn't find configuration in template")
		return nil
	}
	if err = root.Parse(); err != nil {
		utils.FatalPrintln("Couldn't parse root configuration:", err)
	}
	if root.Description != "" {
		utils.OkPrintf("[green]%s[/] - [yellow]%s[/] - [cyan]%s[/]", root.Name, root.Version, root.Description)
	} else {
		utils.OkPrintf("[green]%s[/] - [yellow]%s[/]", root.Name, root.Version)
	}
	root.Prompt(builtins)
	return root
}

// collect recursively prompts all variables for the template at dir and all
// its includes, and accumulates a flat list of files to render.
// ctx is the pre-fill context from a parent template (empty at top level).
// Returns the file list and the unified render context (all prompted values
// from this template and all its includes merged together).
func collect(dir, output string, ctx conf.InputCtx, depth int) ([]*conf.File, map[string]any, error) {
	var err error

	builtins := conf.DefaultBuiltins(output, nil)
	root := HandleRootConfig(dir, ctx, builtins)
	builtins["template_name"] = root.Name
	builtins["template_version"] = root.Version
	var candidates []*conf.File

	m := make(map[string]*conf.ConfigFile)
	m[root.File.Dir] = &root.ConfigFile

	// Cycle through to find override configuration files
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == ConfigName && path != root.File.Path {
			cf := conf.NewConfigFile(path, info, ctx)
			m[cf.File.Dir] = cf
			utils.OkPrintf("Override configuration: [yellow]%s[/]", path)
			if err := cf.Parse(); err != nil {
				return fmt.Errorf("could not parse configuration: %w", err)
			}
			cf.Prompt(builtins)
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not read filesystem: %w", err)
	}

	// Cycle through the files and attach their renderers
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() != ConfigName && info.Name() != ".git" {
			f := conf.NewFile(path, info, ctx)
			f.Builtins = builtins
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
		return nil, nil, fmt.Errorf("could not read filesystem: %w", err)
	}

	// Build the accumulated render context: start with the incoming pre-fill
	// values (from parent), then overlay built-ins as lowest priority,
	// then overlay this template's prompted results.
	accumCtx := conf.InputCtxToMap(ctx)
	maps.Copy(accumCtx, builtins)
	for _, cf := range m {
		if cf.Variables != nil {
			cf.Variables.AddToCtx("", accumCtx)
		}
	}

	// Process includes: evaluate gates, fetch, recurse.
	for _, inc := range root.Includes {
		if inc.Confirm != nil {
			msg := inc.Prompt
			if msg == "" {
				msg = "Include " + inc.Source + "?"
			}
			var confirmed bool
			if err = survey.AskOne(&survey.Confirm{Message: msg, Default: *inc.Confirm}, &confirmed, nil); err != nil {
				return nil, nil, fmt.Errorf("could not get confirmation for include %q: %w", inc.Source, err)
			}
			if !confirmed {
				utils.OkPrintf("Skipped include [green]%s[/]", inc.Source)
				continue
			}
		}
		if inc.If != "" {
			pass, condErr := conf.EvalCondition(inc.If, accumCtx)
			if condErr != nil {
				utils.ErrPrintf("Condition error for include [green]%s[/] - [red]%s[/]", inc.Source, condErr.Error())
				continue
			}
			if !pass {
				utils.OkPrintf("Skipped include [green]%s[/]", inc.Source)
				continue
			}
		}

		effectiveOutput := output
		if inc.Dest != "" && inc.Dest != "." {
			effectiveOutput = filepath.Join(output, inc.Dest)
		}

		// Resolve relative local paths against the parent template directory.
		source := inc.Source
		if !strings.HasSuffix(source, ".git") && !filepath.IsAbs(source) {
			source = filepath.Join(dir, source)
		}

		p := provider.NewProviderFromPath(source, inc.Path, "", depth)
		if p.UsesTmp() {
			var incTempDir string
			if incTempDir, err = os.MkdirTemp("", "qk-include"); err != nil {
				return nil, nil, fmt.Errorf("could not create temp dir for include %q: %w", source, err)
			}
			defer os.RemoveAll(incTempDir) //nolint:errcheck
			p = provider.NewProviderFromPath(source, inc.Path, incTempDir, depth)
		}

		utils.DebugPrintf("Fetching include via [green]%s[/] provider: [green]%s[/]", p.Name(), source)
		var tpath string
		if tpath, err = p.Fetch(); err != nil {
			return nil, nil, fmt.Errorf("could not fetch include %q: %w", source, err)
		}

		// Recurse: pass the current accumCtx so the include can reuse already-
		// prompted variables without re-prompting.
		incFiles, incCtx, err := collect(tpath, effectiveOutput, conf.MapToInputCtx(accumCtx), depth)
		if err != nil {
			return nil, nil, fmt.Errorf("could not collect include %q: %w", source, err)
		}
		// Merge new variables introduced by the include (parent takes priority
		// on any name collision — but in practice FillPrompt ensures they agree).
		for k, v := range incCtx {
			if _, exists := accumCtx[k]; !exists {
				accumCtx[k] = v
			}
		}
		candidates = append(candidates, incFiles...)
	}

	return candidates, accumCtx, nil
}

// Analyze is the main entry point for rendering a template.
// It runs in two phases:
//  1. Collect: prompt all variables across the full template tree (parent +
//     all includes, recursively) before writing any files.
//  2. Render: write all collected files using the unified context, so that
//     every file has access to variables from every template in the tree.
//
// parentCtx contains values already collected by a parent template; pass an
// empty InputCtx at the top level.
func Analyze(dir, output, input string, set []string, depth int, parentCtx conf.InputCtx) error {
	var err error
	ctx := parentCtx

	if input != "" {
		inputCtx, err := conf.GetInputContext(input)
		if err != nil {
			return fmt.Errorf("could not parse input file: %w", err)
		}
		ctx = conf.MergeCtx(ctx, inputCtx)
		utils.OkPrintf("Input file [green]%s[/] found", input)
	}
	if len(set) > 0 {
		setCtx, err := conf.GetSetContext(set)
		if err != nil {
			return fmt.Errorf("could not parse set flags: %w", err)
		}
		ctx = conf.MergeCtx(ctx, setCtx)
		utils.OkPrintln("Command line set merged in context")
	}

	// Phase 1: prompt all variables across the full template tree.
	candidates, globalCtx, err := collect(dir, output, ctx, depth)
	if err != nil {
		return err
	}

	// Inject the unified context into every file so that variables from
	// included sub-templates are available when rendering parent files
	// (e.g. .license from a license include available in README.md).
	globalInputCtx := conf.MapToInputCtx(globalCtx)
	for _, f := range candidates {
		f.GlobalCtx = globalCtx
		// Also update Ctx so per-file frontmatter prompts can be pre-filled
		// with the full global context.
		f.Ctx = conf.MergeCtx(f.Ctx, globalInputCtx)
	}

	// Phase 2: render all files.
	for _, f := range candidates {
		if err = f.ParseFrontMatter(); err != nil {
			return fmt.Errorf("could not parse front matter for file %s: %w", f.Path, err)
		}
		if err = f.Render(); err != nil {
			return fmt.Errorf("could not render template: %w", err)
		}
	}

	return nil
}
