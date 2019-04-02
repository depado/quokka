package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/Depado/projectmpl/cmd"
	"github.com/Depado/projectmpl/renderer"
	"github.com/spf13/cobra"
)

// Build number and versions injected at compile time
var (
	Version = "unknown"
	Build   = "unknown"
)

var qkdesc = `Quokka (qk) is a template engine that enables to render local or distant 
templates/boilerplates in a user friendly way. When given a URL/Git repository
or a path to a local Quokka template, quokka will ask for the required values
in an interactive way except if an inpute file is given to the CLI.
`

// Main command that will be run when no other command is provided on the
// command-line
var rootc = &cobra.Command{
	Use:   "qk [template] [output] <options>",
	Short: "qk is a boilerplate engine",
	Long:  qkdesc,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		renderer.Render(
			args[0],
			args[1],
			viper.GetString("output"),
			viper.GetString("path"),
			viper.GetString("input"),
			viper.GetBool("keep"),
			viper.GetInt("git.depth"),
		)
	},
}

// Version command that will display the build number and version (if any)
var versionc = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Printf("Build: %s\nVersion: %s\n", Build, Version) },
}

// New command that will create a new empty quokka template
var newc = &cobra.Command{
	Use:   "new [output] <options>",
	Short: "Create a new quokka template",
	Args:  cobra.MinimumNArgs(1),
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("Coming soon") },
}

func main() {
	// Initialize Cobra and Viper
	cobra.OnInitialize(cmd.Initialize)
	cmd.AddRendererFlags(rootc)
	rootc.AddCommand(versionc)
	rootc.AddCommand(newc)

	// Run the command
	if err := rootc.Execute(); err != nil {
		log.Fatal("Couldn't start program:", err)
	}
}
