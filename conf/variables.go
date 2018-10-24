package conf

import (
	"fmt"

	"github.com/Depado/projectmpl/utils"
	"github.com/fatih/color"
	survey "gopkg.in/AlecAivazis/survey.v1"
	yaml "gopkg.in/yaml.v2"
)

// Variables represents a map of variable
type Variables struct {
	m map[string]*Variable
	s []*Variable
}

// FromMapSlice fills in the Variables struct with the data stored in a
// yaml.MapSlice. Used to recursively parse nested variables.
func (vv *Variables) FromMapSlice(in yaml.MapSlice) {
	for _, i := range in {
		inv := &Variable{}
		inv.FromMapItem(i)

		k := i.Key.(string)
		inv.Name = k
		vv.m[k] = inv
		vv.s = append(vv.s, inv)
	}
}

// UnmarshalYAML defines a custom way to unmarshal to the Variables type.
// Specifically this allows to conserve the key order
func (vv *Variables) UnmarshalYAML(unmarshal func(interface{}) error) error {
	variables := Variables{
		m: make(map[string]*Variable),
	}
	n := yaml.MapSlice{}
	err := unmarshal(&n)
	if err != nil {
		return err
	}
	variables.FromMapSlice(n)
	*vv = variables
	return nil
}

// Prompt will prompt
func (vv Variables) Prompt() {
	for _, v := range vv.s {
		v.Prompt()
	}
}

// Ctx generates the context from the variables
func (vv Variables) Ctx() map[string]interface{} {
	ctx := make(map[string]interface{})
	for _, v := range vv.s {
		if v != nil {
			if v.Confirm != nil {
				ctx[v.Name] = *v.Confirm
			} else {
				ctx[v.Name] = v.Result
			}
		}
		if v.Variables != nil {
			v.Variables.AddToCtx(v.Name, ctx)
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
	Confirm   *bool      `yaml:"confirm,omitempty"`
	Variables *Variables `yaml:"variables,omitempty"`

	Result string
	Name   string
}

// FromMapItem will fill the variable with the data stored in the input
// yaml.MapItem. Used to recursively parse nested variables.
func (v *Variable) FromMapItem(i yaml.MapItem) {
	for _, data := range i.Value.(yaml.MapSlice) {
		switch data.Key.(string) {
		case "default":
			v.Default = data.Value.(string)
		case "prompt":
			v.CustomPrompt = data.Value.(string)
		case "values":
			for _, p := range data.Value.([]interface{}) {
				v.Values = append(v.Values, p.(string))
			}
		case "help":
			v.Help = data.Value.(string)
		case "required":
			v.Required = data.Value.(bool)
		case "confirm":
			b := data.Value.(bool)
			v.Confirm = &b
		case "variables":
			vv := &Variables{
				m: make(map[string]*Variable),
			}
			vv.FromMapSlice(data.Value.(yaml.MapSlice))
			v.Variables = vv
		}
	}
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
