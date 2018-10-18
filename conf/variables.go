package conf

import (
	"fmt"

	"github.com/fatih/color"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/Depado/projectmpl/utils"
)

// Variable represents a single variable
type Variable struct {
	// Default value to display to the user for input prompts
	Default string `yaml:"default"`

	// Prompt allows to override the standard prompt and display more info
	CustomPrompt string `yaml:"prompt"`

	// List of possible values for the user to answer
	Values []string `yaml:"values"`

	// If this field isn't empty, then an help message can be shown to the user
	Help string `yaml:"help"`

	// Flags a variable as required, preventing the user from entering empty
	// values
	Required bool `yaml:"required"`

	// Confirm is used both for default variable and to store the result.
	// If this field isn't nil, then a confirmation survey is used.
	Confirm *bool `yaml:"confirm,omitempty"`

	Result string
	Name   string
}

// Prompt prompts for the variable
func (v *Variable) Prompt() {
	var prompt survey.Prompt
	var validator survey.Validator
	msg := fmt.Sprintf("Choose a value for %s:", color.YellowString(v.Name))
	if v.CustomPrompt != "" {
		msg = v.CustomPrompt
	}

	if len(v.Values) != 0 {
		prompt = &survey.Select{
			Message: msg,
			Options: v.Values,
			Default: v.Default,
			Help:    v.Help,
		}
	} else if v.Confirm != nil {
		prompt = &survey.Confirm{
			Message: msg,
			Default: *v.Confirm,
			Help:    v.Help,
		}
	} else {
		prompt = &survey.Input{
			Message: msg,
			Default: v.Default,
			Help:    v.Help,
		}
	}

	if v.Required {
		validator = survey.Required
	}
	var out interface{}
	if v.Confirm != nil {
		out = v.Confirm
	} else {
		out = &v.Result
	}
	if err := survey.AskOne(prompt, out, validator); err != nil {
		utils.FatalPrintln("Couldn't get an answer:", err)
	}
}
