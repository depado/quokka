package renderer

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Depado/projectmpl/provider"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"gopkg.in/AlecAivazis/survey.v1"
)

// Color stuff
var green = color.New(color.FgGreen)
var prefix = green.Sprint("»")
var errprefix = color.New(color.FgRed).Sprint("»")

// FetchTemplate will determine which provider needs to be used to fetch the
// template and display user-friendly information about what's going on under
// the hood
func FetchTemplate(template string) string {
	var path string
	var err error

	// Setup colors and spinner
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	p := provider.NewProviderFromPath(template)
	fmt.Println(prefix, "Detected", green.Sprint(p.Name()), "template provider")
	s.Suffix = fmt.Sprintf(" %s…", strings.Title(p.Action()))
	s.Color("green") // nolint: errcheck
	s.Start()

	// Actually fetch and fail fast if an error occurs
	path, err = p.Fetch()
	if err != nil {
		s.FinalMSG = fmt.Sprintln(errprefix, "Couldn't complete operation:", err)
		s.Stop()
		os.Exit(1)
	}
	s.FinalMSG = fmt.Sprintln(prefix, "Done", p.Action(), "in", green.Sprint(path))
	s.Stop()

	return path
}

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
			fmt.Println("Canceled operation")
			os.Exit(0)
		}
	}

	path = FetchTemplate(template)

	defer func(p string) {
		os.RemoveAll(path)
		fmt.Println(prefix, "Removed template", green.Sprint(path))
	}(path)
}
