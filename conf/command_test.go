package conf

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Depado/quokka/utils"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

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
		output string
	}{
		{"shouldn't fail", fields{Cmd: ""}, ""},
		{"test echo", fields{Cmd: "echo", Echo: "Doing nothing"}, utils.OkSprintln("Doing nothing")},
		{"test output", fields{Cmd: "echo test", Output: true}, utils.OkSprintln("test")},
		{"test not found and no stop", fields{Cmd: "xxxxxxxx"}, utils.ErrSprintln(`Couldn't execute command, ignoring: exec: "xxxxxxxx": executable file not found in $PATH`)},
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
			if tt.output != "" {
				assert.Equal(t, tt.output, captureOutput(c.Run))
			} else {
				c.Run()
			}
		})
	}
}
