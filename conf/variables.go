package conf

import (
	"fmt"
	"sort"

	"github.com/Depado/projectmpl/utils"
	"github.com/fatih/color"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Variables represents a map of variable
type Variables map[string]*Variable

// func (e *Variables) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	n := yaml.MapSlice{}
// 	err := unmarshal(&n)
// 	if err != nil {
// 		return err
// 	}
// 	for _, v := range n {
// 		var inv &Variable
// 		fmt.Println("============")
// 		fmt.Println(v.Key)
// 		for _, vv := range v.Value.(yaml.MapSlice) {
// 			switch vv.Key {
// 			case "default":

// 			}
// 			fmt.Println("\t", vv.Key, vv.Value)
// 		}
// 	}
// 	return nil
// }

// Prompt will prompt the variables
func (vv Variables) Prompt() {
	// Order the variables alphabetically to keep the same order
	var ordered []*Variable
	for k, v := range vv {
		if v == nil { // Unconfigured values do have a key but no value
			v = &Variable{Name: k}
		} else {
			v.Name = k
		}
		ordered = append(ordered, v)
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Name < ordered[j].Name
	})

	for _, variable := range ordered {
		variable.Prompt()
	}
}

// Ctx generates the context from the variables
func (vv Variables) Ctx() map[string]interface{} {
	ctx := make(map[string]interface{})
	for k, v := range vv {
		if v != nil {
			if v.Confirm != nil {
				ctx[k] = *v.Confirm
			} else {
				ctx[k] = v.Result
			}
		}
		if v.Variables != nil {
			v.Variables.AddToCtx(k, ctx)
		}
	}
	return ctx
}

// AddToCtx will add the variable results to a sub-key
func (vv Variables) AddToCtx(prefix string, ctx map[string]interface{}) {
	for k, v := range vv.Ctx() {
		if prefix != "" {
			ctx[prefix+"_"+k] = v
		} else {
			ctx[k] = v
		}
	}
}

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
	Confirm   *bool     `yaml:"confirm,omitempty"`
	Variables Variables `yaml:"variables,omitempty"`

	Result string
	Name   string
}

// True returns if the variable has been filled
func (v *Variable) True() bool {
	return v.Result != "" || v.Confirm != nil && *v.Confirm
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
	if v.True() && v.Variables != nil {
		v.Variables.Prompt()
	}
}
