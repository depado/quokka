package conf

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/expr-lang/expr"
	"github.com/fatih/color"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	"github.com/depado/quokka/utils"
)

const frontMatterPrefix = "---"

// File represents a single file, combining both its path and its os.FileInfo
type File struct {
	Path      string
	Dir       string
	NewPath   string
	Info      os.FileInfo
	Renderers []*ConfigFile
	Metadata  *FileMetadata
	Ctx       InputCtx
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
	"uc": unicode.ToUpper,
}

// AddRenderer adds a renderer to a file
func (f *File) AddRenderer(c *ConfigFile) {
	f.Renderers = append(f.Renderers, c)
}

// ParseFrontMatter will parse the front matter and add a renderer to the file
// if needed
func (f *File) ParseFrontMatter() error {
	var err error
	var fd *os.File

	// Open the file
	if fd, err = os.Open(f.Path); err != nil {
		utils.FatalPrintln("Couldn't open candidate:", err)
	}
	defer fd.Close() //nolint:errcheck

	// Scan it and check if there are known delimiters or an end of file
	scanner := bufio.NewScanner(fd)
	if !scanner.Scan() {
		return nil
	}

	// Detected from matter
	if scanner.Text() != frontMatterPrefix {
		return nil
	}

	// Detected from matter
	var found bool
	for scanner.Scan() {
		if scanner.Text() == frontMatterPrefix {
			found = true
		}
	}
	if !found {
		return nil
	}

	if _, err = fd.Seek(0, 0); err != nil {
		return err
	}
	scanner = bufio.NewScanner(fd)

	var in string
	scanner.Scan() // First line, we know it's front matter
	for scanner.Scan() && scanner.Text() != frontMatterPrefix {
		in += scanner.Text() + "\n"
	}

	// Parse stuff to configuration
	var fm FileMetadata
	if err = yaml.Unmarshal([]byte(in), &fm); err != nil {
		return err
	}
	f.Metadata = &fm

	if f.Metadata.Variables != nil && len(*f.Metadata.Variables) > 0 {
		utils.OkPrintln("Variables for single file", color.YellowString(f.Path))
		f.Metadata.Variables.FillPrompt("", f.Ctx)
	}
	return nil
}

// WriteCopy will write the file to its intended path and not attempt to
// render
func (f *File) WriteCopy() error {
	var err error
	var ofd *os.File // Output
	var sfd *os.File // Source

	// Create the directory
	if err = os.MkdirAll(filepath.Dir(f.NewPath), os.ModePerm); err != nil {
		return err
	}

	if sfd, err = os.Open(f.Path); err != nil {
		return err
	}
	defer sfd.Close() //nolint:errcheck
	if ofd, err = os.Create(f.NewPath); err != nil {
		return err
	}
	defer ofd.Close() //nolint:errcheck

	// Scan it and check if there are known delimiters or an end of file
	scanner := bufio.NewScanner(sfd)
	if !scanner.Scan() {
		return nil
	}

	// Detected from matter
	var found bool
	if scanner.Text() == frontMatterPrefix {
		for scanner.Scan() {
			if scanner.Text() == frontMatterPrefix {
				found = true
				break
			}
		}
		if !found {
			if _, err = sfd.Seek(0, 0); err != nil {
				return err
			}
			scanner = bufio.NewScanner(sfd)
		}
	} else {
		if _, err = ofd.WriteString(scanner.Text() + "\n"); err != nil {
			return err
		}
	}
	for scanner.Scan() {
		if _, err = ofd.WriteString(scanner.Text() + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// WriteRender will first render the file as if ignored, but will parse it and
// render it as soon as it has been copied
func (f *File) WriteRender(ctx map[string]interface{}, delims []string) error {
	var err error
	var fd *os.File
	rdr := f.NewPath + ".rendered"

	if err = f.WriteCopy(); err != nil {
		return err
	}

	t := template.Must(template.New(path.Base(f.NewPath)).Delims(delims[0], delims[1]).Funcs(templateFuncMaps).ParseFiles(f.NewPath))

	if fd, err = os.Create(f.NewPath + ".rendered"); err != nil {
		return err
	}

	if err = t.Execute(fd, ctx); err != nil {
		fd.Close() //nolint:errcheck
		return err
	}

	fd.Close() //nolint:errcheck
	return os.Rename(rdr, f.NewPath)
}

// Render will actually render the file
func (f *File) Render() error {
	var err error
	var condition string
	var copy bool
	var ignore bool

	ctx := make(map[string]interface{})
	delims := []string{"{{", "}}"}
	for i := len(f.Renderers) - 1; i >= 0; i-- {
		r := f.Renderers[i]
		if r.Copy != nil {
			copy = *r.Copy
		}
		if r.Ignore != nil {
			ignore = *r.Ignore
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
			copy = *f.Metadata.Copy
		}
		if f.Metadata.Ignore != nil {
			ignore = *f.Metadata.Ignore
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
	if ignore {
		utils.OkPrintln("Ignored ", color.GreenString(f.NewPath))
		return nil
	}
	if condition != "" {
		if len(strings.Fields(condition)) == 1 {
			if v, ok := ctx[condition]; ok {
				switch o := v.(type) {
				case bool:
					if !o {
						utils.OkPrintln("Ignored ", color.GreenString(f.NewPath))
						return nil
					}
				case string:
					if o == "" {
						utils.OkPrintln("Ignored ", color.GreenString(f.NewPath))
						return nil
					}
				}
			}
		} else {
			p, err := expr.Compile(condition, expr.AsBool())
			if err != nil {
				utils.ErrPrintln("Invalid conditional in", color.YellowString(f.Path), "-", color.RedString(err.Error()))
				return nil
			}
			out, err := expr.Run(p, ctx)
			if err != nil {
				utils.ErrPrintln("Failed to run condition in", color.YellowString(f.Path), "-", color.RedString(err.Error()))
				return nil
			}
			res, ok := out.(bool)
			if !ok {
				utils.ErrPrintln("Condition didn't return a boolean value in", color.YellowString(f.Path))
				return nil
			}
			if !res {
				utils.OkPrintln("Ignored ", color.GreenString(f.NewPath))
				return nil
			}
		}
	}
	if copy {
		if err = f.WriteCopy(); err != nil {
			return err
		}
		utils.OkPrintln("Copied  ", color.GreenString(f.NewPath))
	} else {
		if err = f.WriteRender(ctx, delims); err != nil {
			return err
		}
		utils.OkPrintln("Rendered", color.GreenString(f.NewPath))
	}
	return nil
}
