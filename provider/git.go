package provider

import (
	"io/ioutil"

	"github.com/spf13/viper"
	git "gopkg.in/src-d/go-git.v4"
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

func (gitp) Action() string {
	return "cloning"
}

func (g gitp) Fetch() (string, error) {
	var err error
	var tmpdir string

	if tmpdir, err = ioutil.TempDir("", "projectmpl"); err != nil {
		return "", err
	}
	_, err = git.PlainClone(tmpdir, false, &git.CloneOptions{
		Depth: viper.GetInt("git.depth"),
		URL:   g.Repo,
	})
	if err != nil {
		return "", err
	}
	return tmpdir, nil
}
