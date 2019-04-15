package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/AlecAivazis/survey.v1"

	_ "github.com/Depado/quokka/conf" // Import conf to ensure package init for survey
	"github.com/Depado/quokka/utils"
)

// NewQuokkaTemplate will create a new Quokka template with default params
func NewQuokkaTemplate(path string) {
	var err error

	if _, err = os.Stat(path); !os.IsNotExist(err) {
		var confirmed bool
		prompt := &survey.Confirm{
			Message: "The output destination already exists. Continue ?",
		}
		survey.AskOne(prompt, &confirmed, nil) // nolint: errcheck
		if !confirmed {
			utils.ErrPrintln("Canceled operation")
			os.Exit(0)
		}
	} else {
		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			utils.FatalPrintln("Unable to create directory")
		}
	}
	var qs = []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Name of the template?",
				Default: "Quokka Template",
			},
		},
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Description of the template?",
				Default: "New Quokka Template",
			},
		},
		{
			Name: "version",
			Prompt: &survey.Input{
				Message: "Version of the template?",
				Default: "0.1.0",
			},
		},
	}
	answers := struct {
		Name        string // survey will match the question and field names
		Description string `survey:"color"` // or you can tag fields to match a specific name
		Version     int    // if the types don't match exactly, survey will try to convert for you
	}{}

	// perform the questions
	if err = survey.Ask(qs, &answers); err != nil {
		utils.FatalPrintln("Unable to ask questions")
	}

	fmt.Println(answers)
}
