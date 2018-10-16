package renderer

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/Depado/projectmpl/provider"
	"github.com/Depado/projectmpl/utils"
)

// Render is the main render function
func Render(template, output string) {
	var err error
	var path string

	if _, err = os.Stat(output); !os.IsNotExist(err) {
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

	// Determines the provider to use and fetch the template
	p := provider.NewProviderFromPath(template)
	utils.OkPrintln("Detected", utils.Green.Sprint(p.Name()), "template provider")
	if path, err = p.Fetch(); err != nil {
		os.Exit(1)
	}

	// Delete the template if needed
	if !viper.GetBool("keep") && p.UsesTmp() {
		defer func(p string) {
			os.RemoveAll(path)
			utils.OkPrintln("Removed template", utils.Green.Sprint(path))
		}(path)
	}
	Analyze(path)
}
