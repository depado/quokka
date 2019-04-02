package utils

import (
	"io/ioutil"
)

// GetTemplateDir creates the necessary directory or uses the one provided in
// the configuration on the command line
func GetTemplateDir(output string) (string, error) {
	if output != "" {
		return output, nil
	}
	return ioutil.TempDir("", "qk")
}
