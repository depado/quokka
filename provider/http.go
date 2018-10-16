package provider

import (
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/mholt/archiver"

	"github.com/Depado/projectmpl/utils"
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

func (h httpp) Unarchive(fn, outdir string, from io.Reader) error {
	if extractor := archiver.MatchingFormat(fn); extractor != nil {
		return extractor.Read(from, outdir)
	}
	return fmt.Errorf("not an archive")
}

func (h httpp) Fetch() (string, error) {
	var err error
	var resp *http.Response
	var outdir string

	// Setup utils and spinner
	s := utils.NewSpinner("Downloading template")
	// Create directory if needed
	if outdir, err = utils.GetTemplateDir(); err != nil {
		s.ErrStop("Couldn't create template directory:", err)
		return "", err
	}
	// Download the file using GET method
	if resp, err = http.Get(h.URL); err != nil {
		s.ErrStop("Couldn't download file:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Extract the file
	s.Suffix = " Extracting template…"
	if err = h.Unarchive(path.Base(h.URL), outdir, resp.Body); err != nil {
		s.ErrStop("Couldn't extract archive:", err)
		return "", err
	}
	s.DoneStop("Donwloaded and extracted template in", utils.Green.Sprint(outdir))

	return outdir, nil
}

func (httpp) Name() string {
	return "http"
}

func (httpp) UsesTmp() bool {
	return true
}
