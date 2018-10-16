package provider

type local struct {
	Path string
}

// NewLocalProvider will return a new provider to a local filesystem path
func NewLocalProvider(in string) Provider {
	return local{
		Path: in,
	}
}

func (l local) Fetch() (string, error) {
	return l.Path, nil
}

func (local) UsesTmp() bool {
	return false
}

func (local) Name() string {
	return "local"
}
