package conf

import (
	"os/exec"
	"strings"

	"github.com/Depado/quokka/utils"
)

// Command is a simple command to be executed
type Command struct {
	Cmd     string `yaml:"cmd"`
	Failure string `yaml:"failure"`
	Output  bool   `yaml:"output"`
	Echo    string `yaml:"echo"`
	If      string `yaml:"if"`
}

// Parse wil parse the command line and return an exec.Cmd. If the command is
// empty, nil is returned
func (c Command) Parse() *exec.Cmd {
	var main string
	var args []string

	parts := strings.Fields(c.Cmd)
	if len(parts) == 0 {
		return nil
	} else if len(parts) == 1 {
		main = parts[0]
	} else if len(parts) > 1 {
		main = parts[0]
		args = parts[1:]
	}

	if len(args) > 0 {
		return exec.Command(main, args...)
	}
	return exec.Command(main)
}

// Run will run a single command
func (c Command) Run() {
	var err error
	var output []byte

	cmd := c.Parse()
	if cmd == nil {
		return
	}
	if output, err = cmd.Output(); err != nil {
		if c.Failure == "stop" {
			utils.FatalPrintln("Couldn't run after command:", err)
		}
		utils.ErrPrintln("Couldn't execute command, ignoring:", err)
		return
	}
	if c.Output {
		out := strings.Split(string(output), "\n")
		for i := 0; i < len(out)-1; i++ {
			utils.OkPrintln(out[i])
		}
	}
	if c.Echo != "" {
		utils.OkPrintln(c.Echo)
	}
}
