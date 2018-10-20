package conf

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestCommand_Parse(t *testing.T) {
	tests := []struct {
		name string
		cmd  string
		want *exec.Cmd
	}{
		{"should be nil", "", nil},
		{"match simple command", "ls", exec.Command("ls")},
		{"match one arg", "ls -l", exec.Command("ls", "-l")},
		{"match multiple args", "ls -l -a", exec.Command("ls", "-l", "-a")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Command{
				Cmd: tt.cmd,
			}
			if got := c.Parse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_Run(t *testing.T) {
	type fields struct {
		Cmd     string
		Failure string
		Output  bool
		Echo    string
		If      string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"shouldn't fail", fields{Cmd: ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Command{
				Cmd:     tt.fields.Cmd,
				Failure: tt.fields.Failure,
				Output:  tt.fields.Output,
				Echo:    tt.fields.Echo,
				If:      tt.fields.If,
			}
			c.Run()
		})
	}
}
