package cmd

import (
	"fmt"
	"os"

	"gopkg.in/AlecAivazis/survey.v1"

	_ "github.com/Depado/quokka/conf" // Import conf to ensure package init for survey
	"github.com/Depado/quokka/utils"
)

func askIfNotString(in *string, name, message, def string, debug bool) {
	var err error
	if *in == "" {
		if err = survey.AskOne(&survey.Input{
			Message: message,
			Default: def,
		}, in, nil); err != nil {
			utils.ErrPrintln("Canceled operation")
			os.Exit(0)
		}
	} else if debug {
		utils.OkPrintln(utils.Green.Sprint(name), "already filled:", *in)
	}
}

// NewQuokkaTemplate will create a new Quokka template with default params
func NewQuokkaTemplate(path, name, description, version string, yes, debug bool) {
	var err error

	if _, err = os.Stat(path); !os.IsNotExist(err) {
		if yes {
			utils.OkPrintln("Output destination already exists but 'yes' option was used")
		} else {
			var confirmed bool
			prompt := &survey.Confirm{
				Message: "The output destination already exists. Continue ?",
			}
			survey.AskOne(prompt, &confirmed, nil) // nolint: errcheck
			if !confirmed {
				utils.ErrPrintln("Canceled operation")
				os.Exit(0)
			}
		}
	} else {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			utils.FatalPrintln("Unable to create directory")
		}
	}

	askIfNotString(&name, "name", "Template name?", "Quokka Template", debug)
	askIfNotString(&description, "description", "Template description?", "New Quokka Template", debug)
	askIfNotString(&version, "version", "Template version?", "0.1.0", debug)

	fmt.Printf("Name: %s\nDescription: %s\nVersion: %s\n", name, description, version)
}
