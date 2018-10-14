package provider

type git struct {
}

// NewGitProvider will return a new provider from a git repository
func NewGitProvider(url string) Provider {
	return git{}
}

func (git) Fetch(string) error {
	panic("not implemented")
}
