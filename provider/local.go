package provider

type local struct {
}

// NewLocalProvider will return a new provider to a local filesystem path
func NewLocalProvider(in string) Provider {
	return local{}
}

func (local) Fetch(string) error {
	panic("not implemented")
}
