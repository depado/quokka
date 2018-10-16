package provider

import (
	"path/filepath"

	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"

	"github.com/Depado/projectmpl/utils"
)

type gitp struct {
	Repo string
	Path string
}

// NewGitProvider will return a new provider from a git repository
func NewGitProvider(url string) Provider {
	return &gitp{Repo: url}
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
	if outdir, err = utils.GetTemplateDir(); err != nil {
		s.ErrStop("Couldn't create template directory:", err)
		return "", err
	}
	// Clone the given repository
	if _, err = git.PlainClone(outdir, false, &git.CloneOptions{Depth: viper.GetInt("git.depth"), URL: g.Repo}); err != nil {
		s.ErrStop("Couldn't clone repo:", err)
		return "", err
	}
	s.DoneStop("Done cloning in", utils.Green.Sprint(outdir))

	return filepath.Join(outdir, viper.GetString("template.path")), nil
}

func (gitp) UsesTmp() bool {
	return true
}
