package conf

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/Depado/projectmpl/utils"
	"gopkg.in/yaml.v2"
)

// AllCandidates is the full list of candidates
var AllCandidates []*File

// File represents a single file, combining both its path and its os.FileInfo
type File struct {
	Path      string
	Dir       string
	NewPath   string
	Info      os.FileInfo
	Renderers []*ConfigFile
	Metadata  *Config
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
	defer fd.Close()

	// Scan it and check if there are known delimiters or an end of file
	scanner := bufio.NewScanner(fd)
	if !scanner.Scan() {
		return nil
	}

	// Detected from matter
	if scanner.Text() != "---" {
		return nil
	}

	// Detected from matter
	var found bool
	for scanner.Scan() {
		if scanner.Text() == "---" {
			found = true
		}
	}
	if !found {
		return nil
	}

	fd.Seek(0, 0)
	scanner = bufio.NewScanner(fd)

	var in string
	scanner.Scan() // First line, we know it's front matter
	for scanner.Scan() && scanner.Text() != "---" {
		in += scanner.Text() + "\n"
	}

	// Parse stuff to configuration
	var r Config
	if err = yaml.Unmarshal([]byte(in), &r); err != nil {
		return err
	}
	f.Metadata = &r
	return nil
}

// WriteIgnore will write the file to its intended path and not attempt to
// render
func (f *File) WriteIgnore() error {
	var err error
	var ofd *os.File // Output
	var sfd *os.File // Source

	if sfd, err = os.Open(f.Path); err != nil {
		return err
	}
	defer sfd.Close()
	if ofd, err = os.Create(f.NewPath); err != nil {
		return err
	}
	defer ofd.Close()

	// Scan it and check if there are known delimiters or an end of file
	scanner := bufio.NewScanner(sfd)
	if !scanner.Scan() {
		return nil
	}

	// Detected from matter
	var found bool
	if scanner.Text() == "---" {
		for scanner.Scan() {
			if scanner.Text() == "---" {
				found = true
				break
			}
		}
		if !found {
			sfd.Seek(0, 0)
			scanner = bufio.NewScanner(sfd)
		}
	} else {
		ofd.WriteString(scanner.Text() + "\n")
	}
	for scanner.Scan() {
		ofd.WriteString(scanner.Text() + "\n")
	}
	return nil
}

// WriteRender will first render the file as if ignored, but will parse it and
// render it as soon as it has been copied
func (f *File) WriteRender(ctx map[string]interface{}, delims []string) error {
	var err error
	var fd *os.File
	rdr := f.NewPath + ".rendered"

	if err = f.WriteIgnore(); err != nil {
		return err
	}

	t := template.Must(template.New(path.Base(f.NewPath)).Delims(delims[0], delims[1]).ParseFiles(f.NewPath))

	if fd, err = os.Create(f.NewPath + ".rendered"); err != nil {
		return err
	}
	defer fd.Close()

	if err = t.Execute(fd, ctx); err != nil {
		return err
	}
	return os.Rename(rdr, f.NewPath)
}

// Render will actually render the file
func (f *File) Render() error {
	delims := []string{"{{", "}}"}
	ctx := make(map[string]interface{})

	if f.Metadata != nil {
		if f.Metadata.Ignore {
			return f.WriteIgnore()
		}
	}

	for i := len(f.Renderers) - 1; i >= 0; i-- {
		r := f.Renderers[i]
		if r.Ignore {
			return f.WriteIgnore()
		}
		for k, v := range r.Variables {
			if v != nil {
				if v.Confirm != nil {
					ctx[k] = *v.Confirm
				} else {
					ctx[k] = v.Result
				}
			}
		}
		if r.Delimiters != nil {
			if len(r.Delimiters) != 2 {
				return fmt.Errorf("Delimiters should be an array of two string")
			}
			delims = r.Delimiters
		}
	}
	if f.Metadata != nil && f.Metadata.Delimiters != nil {
		delims = f.Metadata.Delimiters
	}

	return f.WriteRender(ctx, delims)
}
