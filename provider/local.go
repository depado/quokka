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

func (local) Fetch() (string, error) {
	panic("not implemented")
}

func (local) Name() string {
	return "local"
}

func (local) Action() string {
	return "copying"
}
