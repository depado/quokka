package provider

import (
	"reflect"
	"testing"
)

func TestNewProviderFromPath(t *testing.T) {
	type args struct {
		in     string
		path   string
		output string
		depth  int
	}
	tests := []struct {
		name string
		args args
		want Provider
	}{
		{"should detect git", args{in: "git@github.com:Depado/bfchroma.git", depth: 1}, NewGitProvider("git@github.com:Depado/bfchroma.git", "", "", 1)},
		{"should detect git and tweak depth", args{in: "git@github.com:Depado/bfchroma.git", depth: 10}, NewGitProvider("git@github.com:Depado/bfchroma.git", "", "", 10)},
		{"should detect git with http", args{in: "https://github.com/Depado/bfchroma.git", depth: 1}, NewGitProvider("https://github.com/Depado/bfchroma.git", "", "", 1)},
		{"should detect local", args{in: "/tmp/template/"}, NewLocalProvider("/tmp/template/", "")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProviderFromPath(tt.args.in, tt.args.path, tt.args.output, tt.args.depth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProviderFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
