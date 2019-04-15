package renderer

import (
	"os"

	"github.com/Depado/quokka/provider"
	"github.com/Depado/quokka/utils"
	"gopkg.in/AlecAivazis/survey.v1"
)

// Render is the main render function
func Render(template, output, toutput, path, input string, set []string, keep bool, depth int, yes bool) {
	var err error
	var tpath string

	if _, err = os.Stat(output); !os.IsNotExist(err) {
		if yes {
			utils.OkPrintln("Output destination already exists but 'yes' option was used")
		} else {
			var confirmed bool
			prompt := &survey.Confirm{
				Help:    "qk will only affect already existing files that match the template you're trying to render",
				Message: "The output destination already exists. Continue ?",
			}
			survey.AskOne(prompt, &confirmed, nil) // nolint: errcheck
			if !confirmed {
				utils.ErrPrintln("Canceled operation")
				os.Exit(0)
			}
		}
	}

	// Determines the provider to use and fetch the template
	p := provider.NewProviderFromPath(template, path, toutput, depth)
	utils.OkPrintln("Detected", utils.Green.Sprint(p.Name()), "template provider")
	if tpath, err = p.Fetch(); err != nil {
		os.Exit(1)
	}

	// Delete the template if needed
	if !keep && p.UsesTmp() {
		defer func(p string) {
			os.RemoveAll(path)
			utils.OkPrintln("Removed template", utils.Green.Sprint(path))
		}(path)
	}
	Analyze(tpath, output, input, set)
}
