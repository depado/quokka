package provider

type httpp struct {
	URL string
}

// NewHTTPProvider will return a new HTTP provider
func NewHTTPProvider(in string) Provider {
	return httpp{
		URL: in,
	}
}

func (httpp) Fetch() (string, error) {
	panic("not implemented")
}

func (httpp) Name() string {
	return "http"
}

func (httpp) Action() string {
	return "downloading"
}

func (httpp) TemplatePath() string {
	return ""
}
