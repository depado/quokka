package conf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/goccy/go-yaml"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/depado/quokka/utils"
)

// Variables is a slice of pointer to a single variable
type Variables []*Variable

// FindNamed will find a variable by name in the global variables. Returns nil
// if not found
func (vv Variables) FindNamed(s string) *Variable {
	return vv.FindWithParent(nil, s)
}

// FindWithParent tries to find the variable with a prefix
func (vv Variables) FindWithParent(p *Variable, s string) *Variable {
	for _, v := range vv {
		if p != nil {
			if p.Name+"_"+v.Name == s {
				return v
			}
		} else if v.Name == s {
			return v
		}
		if v.Variables != nil {
			if out := v.Variables.FindWithParent(v, s); out != nil {
				return out
			}
		}
	}
	return nil
}

// FromMapSlice fills in the Variables struct with the data stored in a
// yaml.MapSlice. Used to recursively parse nested variables.
func (vv *Variables) FromMapSlice(in yaml.MapSlice) error {
	for _, i := range in {
		inv := &Variable{}
		if err := inv.FromMapItem(i); err != nil {
			return err
		}
		*vv = append(*vv, inv)
	}
	return nil
}

// UnmarshalYAML defines a custom way to unmarshal to the Variables type.
// Specifically this allows to conserve the key order
func (vv *Variables) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var variables Variables
	n := yaml.MapSlice{}
	if err := unmarshal(&n); err != nil {
		return err
	}
	if err := variables.FromMapSlice(n); err != nil {
		return err
	}
	*vv = variables
	return nil
}

// Prompt will prompt
func (vv Variables) Prompt() {
	for _, v := range vv {
		v.Prompt()
	}
}

// Ctx generates the context from the variables
func (vv Variables) Ctx() map[string]interface{} {
	ctx := make(map[string]interface{})
	for _, v := range vv {
		if v != nil {
			if v.Confirm != nil {
				ctx[v.Name] = *v.Confirm
			} else {
				ctx[v.Name] = v.Result
			}
			if v.Variables != nil {
				v.Variables.AddToCtx(v.Name, ctx)
			}
		}
	}
	return ctx
}

func resolveDefault(v *Variable, builtins map[string]interface{}) {
	if !strings.HasPrefix(v.Default, "$") || builtins == nil {
		return
	}
	name := v.Default[1:]
	if val, ok := builtins[name]; ok {
		switch vt := val.(type) {
		case string:
			v.Default = vt
		case int:
			v.Default = strconv.Itoa(vt)
		case int64:
			v.Default = strconv.FormatInt(vt, 10)
		case float64:
			v.Default = strconv.FormatFloat(vt, 'f', -1, 64)
		case bool:
			v.Default = strconv.FormatBool(vt)
		}
	}
}

// FillPrompt will fill the variables from the input context if needed
func (vv *Variables) FillPrompt(prefix string, ctx InputCtx, builtins map[string]interface{}) {
	for _, v := range *vv {
		resolveDefault(v, builtins)

		var ok bool
		p := v.Name
		if prefix != "" {
			p = prefix + "_" + v.Name
		}
		for _, in := range ctx {
			if in.Key == p {
				ok = true
				v.FillFromMapItem(in)
				if v.True() && v.Variables != nil {
					v.Variables.FillPrompt(p, ctx, builtins)
				}
			}
		}
		if !ok {
			v.Prompt()
		}
	}
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

// FillFromMapItem fills the value not from prompt but from a mapitem
func (v *Variable) FillFromMapItem(i yaml.MapItem) {
	if v.Confirm != nil {
		b, ok := i.Value.(bool)
		if !ok {
			utils.ErrPrintln("wrong type for", v.Name, "expecting bool")
			return
		}
		v.Confirm = &b
	} else {
		s, ok := i.Value.(string)
		if !ok {
			utils.ErrPrintln("wrong type for", v.Name, "expecting string")
			return
		}
		v.Result = s
	}
}

// FromMapItem will fill the variable with the data stored in the input
// yaml.MapItem. Used to recursively parse nested variables.
func (v *Variable) FromMapItem(i yaml.MapItem) error {
	key, ok := i.Key.(string)
	if !ok {
		return fmt.Errorf("variable key must be a string, got %T", i.Key)
	}
	v.Name = key

	if i.Value == nil {
		return nil
	}
	fields, ok := i.Value.(yaml.MapSlice)
	if !ok {
		return fmt.Errorf("variable %q: expected a mapping, got %T", v.Name, i.Value)
	}

	for _, data := range fields {
		fkey, ok := data.Key.(string)
		if !ok {
			return fmt.Errorf("variable %q: field key must be a string, got %T", v.Name, data.Key)
		}
		switch fkey {
		case "default":
			s, ok := data.Value.(string)
			if !ok {
				return fmt.Errorf("variable %q: 'default' must be a string, got %T", v.Name, data.Value)
			}
			v.Default = s
		case "prompt":
			s, ok := data.Value.(string)
			if !ok {
				return fmt.Errorf("variable %q: 'prompt' must be a string, got %T", v.Name, data.Value)
			}
			v.CustomPrompt = s
		case "values":
			items, ok := data.Value.([]interface{})
			if !ok {
				return fmt.Errorf("variable %q: 'values' must be a list, got %T", v.Name, data.Value)
			}
			for _, p := range items {
				s, ok := p.(string)
				if !ok {
					return fmt.Errorf("variable %q: 'values' entries must be strings, got %T", v.Name, p)
				}
				v.Values = append(v.Values, s)
			}
		case "help":
			s, ok := data.Value.(string)
			if !ok {
				return fmt.Errorf("variable %q: 'help' must be a string, got %T", v.Name, data.Value)
			}
			v.Help = s
		case "required":
			b, ok := data.Value.(bool)
			if !ok {
				return fmt.Errorf("variable %q: 'required' must be a bool, got %T", v.Name, data.Value)
			}
			v.Required = b
		case "confirm":
			b, ok := data.Value.(bool)
			if !ok {
				return fmt.Errorf("variable %q: 'confirm' must be a bool, got %T", v.Name, data.Value)
			}
			v.Confirm = &b
		case "variables":
			ms, ok := data.Value.(yaml.MapSlice)
			if !ok {
				return fmt.Errorf("variable %q: 'variables' must be a mapping, got %T", v.Name, data.Value)
			}
			var vv Variables
			if err := vv.FromMapSlice(ms); err != nil {
				return err
			}
			v.Variables = vv
		}
	}
	return nil
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
