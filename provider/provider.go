package provider

import (
	"strings"
)

// Provider is the main interface. A provider defines a way to retrieve a
// template and the way said provider should behave once fetched
type Provider interface {
	Fetch() (string, error) // Fetch the template and return the path where it's stored on the local filesystem
	Name() string           // For example, git provider should return "git provider detected"
	UsesTmp() bool          // Defines if the provider uses a temporary file/directory that can be safely deleted at the end
}

// NewProviderFromPath will return a new Provider according to the given string
// as it will try to detect git repositories, http scheme and fallback to local
// directory otherwise
func NewProviderFromPath(in, path, output string, depth int) Provider {
	if strings.HasSuffix(in, ".git") {
		return NewGitProvider(in, path, output, depth)
	}
	if strings.HasPrefix(in, "http://") || strings.HasPrefix(in, "https://") {
		return NewHTTPProvider(in, path, output)
	}
	return NewLocalProvider(in, output)
}
