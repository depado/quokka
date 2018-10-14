package provider

type http struct {
}

// NewHTTPProvider will return a new HTTP provider
func NewHTTPProvider(in string) Provider {
	return http{}
}

func (http) Fetch(string) error {
	panic("not implemented")
}
