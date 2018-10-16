package provider

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"github.com/mholt/archiver"

	"github.com/Depado/projectmpl/colors"
	"github.com/briandowns/spinner"
)

type httpp struct {
	URL string
}

// NewHTTPProvider will return a new HTTP provider
func NewHTTPProvider(in string) Provider {
	return httpp{
		URL: in,
	}
}

func (h httpp) Fetch() (string, error) {
	var err error
	var resp *http.Response
	var dir string
	fn := path.Base(h.URL)

	// Setup colors and spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Dowloading template…"
	s.Color("green") // nolint: errcheck
	s.Start()

	// Create a temporary directory
	if dir, err = ioutil.TempDir("", "projectmpl"); err != nil {
		s.FinalMSG = fmt.Sprintln(colors.ErrPrefix, "Couldn't create tmp dir:", err)
		s.Stop()
		return "", err
	}

	// Download the file using GET method
	if resp, err = http.Get(h.URL); err != nil {
		s.FinalMSG = fmt.Sprintln(colors.ErrPrefix, "Couldn't download template:", err)
		s.Stop()
		return "", err
	}
	defer resp.Body.Close()

	// Extract
	s.Suffix = " Extracting template…"
	extractor := archiver.MatchingFormat(fn)
	if extractor != nil {
		if err = extractor.Read(resp.Body, dir); err != nil {
			s.FinalMSG = fmt.Sprintln(colors.ErrPrefix, "Couldn't extract archive:", err)
			s.Stop()
			return "", err
		}
	} else {
		s.FinalMSG = fmt.Sprintln(colors.ErrPrefix, "Unknown archive format or not an archive")
		s.Stop()
		return "", fmt.Errorf("not an archive")
	}

	s.FinalMSG = fmt.Sprintln(colors.OkPrefix, "Downloaded and extracted template in", colors.Green.Sprint(dir))
	s.Stop()

	return dir, nil
}

func (httpp) Name() string {
	return "http"
}

func (httpp) UsesTmp() bool {
	return true
}
