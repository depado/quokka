package utils

import (
	"io/ioutil"

	"github.com/spf13/viper"
)

// GetTemplateDir creates the necessary directory or uses the one provided in
// the configuration on the command line
func GetTemplateDir() (string, error) {
	tmplpath := viper.GetString("template.output")
	if tmplpath != "" {
		return tmplpath, nil
	}
	return ioutil.TempDir("", "projectmpl")
}
