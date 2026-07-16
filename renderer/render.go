package renderer

import (
	"os"

	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/depado/quokka/conf"
	"github.com/depado/quokka/provider"
	"github.com/depado/quokka/utils"
)

// Render is the main render function
func Render(template, output, toutput, path, input string, set []string, keep bool, depth int, yes, trusted, noCommands bool) {
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
	utils.DebugPrintf("Detected [green]%s[/] template provider", p.Name())
	if tpath, err = p.Fetch(); err != nil {
		os.Exit(1)
	}

	// Delete the template if needed
	if !keep && p.UsesTmp() {
		defer func(tp string) {
			os.RemoveAll(tp) //nolint:errcheck
			utils.OkPrintf("Removed template [green]%s[/]", tp)
		}(tpath)
	}
	if err := Analyze(tpath, output, input, set, depth, conf.InputCtx{}, trusted, noCommands); err != nil {
		utils.FatalPrintln(err)
	}
}
