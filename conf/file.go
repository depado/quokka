package conf

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Depado/projectmpl/utils"
)

// AllCandidates is the full list of candidates
var AllCandidates []*File

// File represents a single file, combining both its path and its os.FileInfo
type File struct {
	Path      string
	Dir       string
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
func (f *File) ParseFrontMatter() {
	var err error
	var fd *os.File

	if fd, err = os.Open(f.Path); err != nil {
		utils.FatalPrintln("Couldn't open candidate:", err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	if !scanner.Scan() {
		return
	}
	// Detected from matter
	if scanner.Text() == "---" {
		var line string
		for scanner.Scan() && scanner.Text() != "---" {
			line = scanner.Text()
			fmt.Println(line)
		}
	}
	return
}
