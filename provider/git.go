package provider

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Depado/projectmpl/colors"
	"github.com/briandowns/spinner"
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

func (g gitp) Fetch() (string, error) {
	var err error
	var outdir string

	// Setup colors and spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Cloning…"
	s.Color("green") // nolint: errcheck
	s.Start()

	if viper.GetString("template.path") != "" {
		outdir = viper.GetString("template.path")
	} else {
		// Create the temporary directory where we'll clone the repo
		if outdir, err = ioutil.TempDir("", "projectmpl"); err != nil {
			s.FinalMSG = fmt.Sprintln(colors.ErrPrefix, "Couldn't create tmp dir:", err)
			s.Stop()
			return "", err
		}
	}

	// Clone the given repository
	if _, err = git.PlainClone(outdir, false, &git.CloneOptions{Depth: viper.GetInt("git.depth"), URL: g.Repo}); err != nil {
		s.FinalMSG = fmt.Sprintln(colors.ErrPrefix, "Couldn't clone repo:", err)
		s.Stop()
		return "", err
	}
	s.FinalMSG = fmt.Sprintln(colors.OkPrefix, "Done cloning in", colors.Green.Sprint(outdir))
	s.Stop()
	return outdir, nil
}

func (gitp) UsesTmp() bool {
	return true
}
