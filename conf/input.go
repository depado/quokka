package conf

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// InputCtx is the input context
type InputCtx yaml.MapSlice

// MergeCtx will merge two InputCtx into a single one
func MergeCtx(a, b InputCtx) InputCtx {
	out := a
	for _, v := range b {
		var found bool
		for i, base := range out {
			if base.Key == v.Key {
				out[i].Value = v.Value
				found = true
			}
		}
		if !found {
			out = append(out, v)
		}
	}

	return out
}

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

// GetSetContext will return the map of string to interface{} that contains the
// set flags passed on the command line parsed
func GetSetContext(set []string) (InputCtx, error) {
	out := InputCtx{}
	for _, s := range set {
		tmp := strings.SplitN(s, "=", 2)
		if len(tmp) != 2 {
			return out, fmt.Errorf("invalid set option: %s", s)
		}
		out = append(out, yaml.MapItem{Key: tmp[0], Value: tmp[1]})
	}
	return out, nil
}
