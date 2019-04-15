package utils

import (
	"fmt"
	"os"

	"gopkg.in/AlecAivazis/survey.v1"
)

// AskIfEmptyString will prompt the user if needed
func AskIfEmptyString(in *string, name, message, def string, debug bool) {
	var err error
	if *in == "" {
		if err = survey.AskOne(&survey.Input{Message: message, Default: def}, in, nil); err != nil {
			ExitPrintln("Canceled operation")
		}
	} else if debug {
		OkPrintln(Green.Sprint(name), "already filled:", *in)
	}
}

// ConfirmFileExists will prompt the user with a confirmation if the
// file/directory already exists
// If the user answers "no", the program exits properly, otherwise the function
// returns true if the file existed and false if it doesn't
func ConfirmFileExists(path string, dir, yes, debug bool) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t := "file"
		p := path
		if dir {
			p = path + "/"
			t = "directory"
		}
		if yes {
			OkPrintln(fmt.Sprintf(
				"The destination %s %s already exists but %s option was used.",
				t,
				Green.Sprint(p),
				Yellow.Sprint("yes"),
			))
			return true
		}
		var confirmed bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf(
				"The destination %s %s already exists. Continue ?",
				t,
				Green.Sprint(p),
			),
		}
		survey.AskOne(prompt, &confirmed, nil) // nolint: errcheck
		if !confirmed {
			ExitPrintln("Canceled operation")
		}
		return true
	}
	return false
}
