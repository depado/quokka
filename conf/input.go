package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// InputCtx is the input context
type InputCtx yaml.MapSlice

// GetInputContext will return a map of string to interface{} that will then
// be used to determine whether or not a value from the root config file
// has already been filled
func GetInputContext(path string) (InputCtx, error) {
	var out InputCtx
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return out, err
	}
	return out, yaml.Unmarshal(input, &out)
}
