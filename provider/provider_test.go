package provider

import (
	"reflect"
	"testing"
)

func TestNewProviderFromPath(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name string
		args args
		want Provider
	}{
		{"should detect git", args{"git@github.com:Depado/bfchroma.git"}, NewGitProvider("git@github.com:Depado/bfchroma.git")},
		{"should detect git with http", args{"https://github.com/Depado/bfchroma.git"}, NewGitProvider("https://github.com/Depado/bfchroma.git")},
		{"should detect http", args{"http://example.com/template/"}, NewHTTPProvider("http://example.com/template/")},
		{"should detect http", args{"https://example.com/template/"}, NewHTTPProvider("https://example.com/template/")},
		{"should detect local", args{"/tmp/template/"}, NewLocalProvider("/tmp/template/")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProviderFromPath(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProviderFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
