package conf

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/expr-lang/expr"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/depado/quokka/utils"
)

const frontMatterPrefix = "---"

// parseFrontMatter reads r and returns the YAML content between the opening
// and closing "---" markers, and the byte offset of the first byte after the
// closing marker. Returns ("", 0, nil) when no valid frontmatter is found.
func parseFrontMatter(r io.Reader) (content string, bodyOffset int64, err error) {
	br := bufio.NewReader(r)

	line, err := br.ReadString('\n')
	if err != nil || strings.TrimRight(line, "\r\n") != frontMatterPrefix {
		return "", 0, nil
	}

	offset := int64(len(line))
	var sb strings.Builder
	for {
		line, err = br.ReadString('\n')
		if err == io.EOF {
			return "", 0, nil // no closing ---
		}
		if err != nil {
			return "", 0, err
		}
		offset += int64(len(line))
		if strings.TrimRight(line, "\r\n") == frontMatterPrefix {
			return sb.String(), offset, nil
		}
		sb.WriteString(line)
	}
}

// File represents a single file, combining both its path and its os.FileInfo
type File struct {
	Path      string
	Dir       string
	NewPath   string
	Info      os.FileInfo
	Renderers []*ConfigFile
	Metadata  *FileMetadata
	Ctx       InputCtx
	// GlobalCtx is the unified render context collected from all templates
	// (parent + all includes). It is the lowest-priority layer: renderer
	// variables and frontmatter always override it.
	GlobalCtx map[string]any
	// Builtins holds the built-in system variables available for default
	// resolution and template rendering (username, hostname, year, etc.)
	Builtins map[string]any
}

type FileMetadata struct {
	Config `yaml:",inline"`
	Rename string `yaml:"rename"`
}

// templateFuncMaps adds just simple strings filter to pretty print string
var templateFuncMaps template.FuncMap = template.FuncMap{
	"title": func(str string) string {
		return cases.Title(language.Und).String(str)
	},
	"uc":   strings.ToUpper,
	"uuid": func() string { return uuid.New().String() },
}

// AddRenderer adds a renderer to a file
func (f *File) AddRenderer(c *ConfigFile) {
	f.Renderers = append(f.Renderers, c)
}

// ParseFrontMatter will parse the front matter and add a renderer to the file
// if needed
func (f *File) ParseFrontMatter() error {
	fd, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer fd.Close() //nolint:errcheck

	content, _, err := parseFrontMatter(fd)
	if err != nil {
		return err
	}
	if content == "" {
		return nil
	}

	var fm FileMetadata
	if err = yaml.Unmarshal([]byte(content), &fm); err != nil {
		return err
	}
	f.Metadata = &fm

	if f.Metadata.Variables != nil && len(*f.Metadata.Variables) > 0 {
		utils.OkPrintf("Variables for single file [yellow]%s[/]", f.Path)
		f.Metadata.Variables.FillPrompt("", f.Ctx, f.Builtins)
	}
	return nil
}

// WriteCopy will write the file to its intended path and not attempt to
// render
func (f *File) WriteCopy() error {
	if err := os.MkdirAll(filepath.Dir(f.NewPath), os.ModePerm); err != nil {
		return err
	}

	sfd, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer sfd.Close() //nolint:errcheck

	_, bodyOffset, err := parseFrontMatter(sfd)
	if err != nil {
		return err
	}
	if _, err = sfd.Seek(bodyOffset, io.SeekStart); err != nil {
		return err
	}

	ofd, err := os.Create(f.NewPath)
	if err != nil {
		return err
	}
	defer ofd.Close() //nolint:errcheck

	_, err = io.Copy(ofd, sfd)
	return err
}

// WriteRender will first render the file as if ignored, but will parse it and
// render it as soon as it has been copied
func (f *File) WriteRender(ctx map[string]any, delims []string) error {
	rdr := f.NewPath + ".rendered"

	if err := f.WriteCopy(); err != nil {
		return err
	}

	t, err := template.New(path.Base(f.NewPath)).Delims(delims[0], delims[1]).Funcs(templateFuncMaps).ParseFiles(f.NewPath)
	if err != nil {
		return err
	}

	fd, err := os.Create(rdr)
	if err != nil {
		return err
	}

	if err = t.Execute(fd, ctx); err != nil {
		fd.Close()     //nolint:errcheck
		os.Remove(rdr) //nolint:errcheck
		return err
	}

	fd.Close() //nolint:errcheck
	return os.Rename(rdr, f.NewPath)
}

// EvalCondition evaluates an if: condition string against the provided context.
// Returns true if the file/include should be processed, false if it should be
// skipped. An unknown single-word variable is treated as false (skip).
// Returns an error only when the expression itself is syntactically invalid.
func EvalCondition(condition string, ctx map[string]any) (bool, error) {
	if len(strings.Fields(condition)) == 1 {
		v, ok := ctx[condition]
		if !ok {
			return false, nil
		}
		switch o := v.(type) {
		case bool:
			return o, nil
		case string:
			return o != "", nil
		}
		return true, nil
	}
	p, err := expr.Compile(condition, expr.AsBool())
	if err != nil {
		return false, fmt.Errorf("invalid conditional %q: %w", condition, err)
	}
	out, err := expr.Run(p, ctx)
	if err != nil {
		return false, fmt.Errorf("failed to run conditional %q: %w", condition, err)
	}
	res, ok := out.(bool)
	if !ok {
		return false, fmt.Errorf("conditional %q did not return a boolean", condition)
	}
	return res, nil
}

// Render will actually render the file
func (f *File) Render() error {
	var err error
	var condition string
	var shouldCopy bool
	var shouldIgnore bool

	ctx := make(map[string]any)
	// Lowest priority: built-in system variables (username, hostname, etc.)
	maps.Copy(ctx, f.Builtins)
	// Next: global context from all templates in the tree.
	maps.Copy(ctx, f.GlobalCtx)
	delims := []string{"{{", "}}"}
	for i := len(f.Renderers) - 1; i >= 0; i-- {
		r := f.Renderers[i]
		if r.Copy != nil {
			shouldCopy = *r.Copy
		}
		if r.Ignore != nil {
			shouldIgnore = *r.Ignore
		}
		if r.Delimiters != nil {
			if len(r.Delimiters) != 2 {
				return fmt.Errorf("delimiters should be an array of two string")
			}
			delims = r.Delimiters
		}
		if r.Variables != nil {
			r.Variables.AddToCtx("", ctx)
		}
	}
	if f.Metadata != nil {
		if f.Metadata.If != "" {
			condition = f.Metadata.If
		}
		if f.Metadata.Copy != nil {
			shouldCopy = *f.Metadata.Copy
		}
		if f.Metadata.Ignore != nil {
			shouldIgnore = *f.Metadata.Ignore
		}
		if f.Metadata.Delimiters != nil {
			delims = f.Metadata.Delimiters
		}
		if f.Metadata.Variables != nil {
			f.Metadata.Variables.AddToCtx("", ctx)
		}
		if f.Metadata.Rename != "" {
			f.NewPath = filepath.Join(filepath.Dir(f.NewPath), f.Metadata.Rename)
		}
	}
	if shouldIgnore {
		utils.OkPrintf("Ignored  [green]%s[/]", f.NewPath)
		return nil
	}
	if condition != "" {
		pass, err := EvalCondition(condition, ctx)
		if err != nil {
			utils.ErrPrintf("Condition error in [yellow]%s[/] - [red]%s[/]", f.Path, err.Error())
			return nil
		}
		if !pass {
			utils.OkPrintf("Ignored  [green]%s[/]", f.NewPath)
			return nil
		}
	}
	if shouldCopy {
		if err = f.WriteCopy(); err != nil {
			return err
		}
		utils.OkPrintf("Copied   [green]%s[/]", f.NewPath)
	} else {
		if err = f.WriteRender(ctx, delims); err != nil {
			return err
		}
		utils.OkPrintf("Rendered [green]%s[/]", f.NewPath)
	}
	return nil
}
