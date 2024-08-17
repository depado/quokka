package provider

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"

	"github.com/depado/quokka/utils"
)

type gitp struct {
	Repo      string
	Path      string
	InnerPath string
	Depth     int
	Output    string
}

// NewGitProvider will return a new provider from a git repository
func NewGitProvider(url, path, output string, depth int) Provider {
	return &gitp{
		Repo:      url,
		InnerPath: path,
		Depth:     depth,
		Output:    output,
	}
}

func (gitp) Name() string {
	return "git"
}

func (g gitp) Fetch() (string, error) {
	var err error
	var outdir string

	// Setup spinner
	s := utils.NewSpinner("Cloning template")
	// Create template directory if needed
	if outdir, err = utils.GetTemplateDir(g.Output); err != nil {
		s.ErrStop("Couldn't create template directory:", err)
		return "", err
	}
	// Clone the given repository
	if _, err = git.PlainClone(outdir, false, &git.CloneOptions{Depth: g.Depth, URL: g.Repo}); err != nil {
		s.ErrStop("Couldn't clone repo:", err)
		return "", err
	}
	s.DoneStop("Done cloning in", utils.Green.Sprint(outdir))

	return filepath.Join(outdir, g.InnerPath), nil
}

func (gitp) UsesTmp() bool {
	return true
}
