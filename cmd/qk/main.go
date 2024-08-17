package main

import (
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/depado/quokka/cmd"
	"github.com/depado/quokka/renderer"
	"github.com/depado/quokka/utils"
)

// Build number and versions injected at compile time
var (
	Version   = "unknown"
	Build     = "unknown"
	BuildDate = "unknown"
)

var qkdesc = `Quokka (qk) is a template engine that enables to render local or
distant templates/boilerplates in a user friendly way. When given a Git
repository or a path to a local Quokka template, quokka will ask for the
required values in an interactive way except if an inpute file is given to the
CLI.
`

// Main command that will be run when no other command is provided on the
// command-line
var rootc = &cobra.Command{
	Use:   "qk [template] [output] <options>",
	Short: "qk is a boilerplate engine",
	Long:  qkdesc,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		utils.OkPrintln("quokka", color.GreenString(Version))
		renderer.Render(
			args[0],
			args[1],
			viper.GetString("output"),
			viper.GetString("path"),
			viper.GetString("input"),
			viper.GetStringSlice("set"),
			viper.GetBool("keep"),
			viper.GetInt("git.depth"),
			viper.GetBool("yes"),
		)
	},
}

// Version command that will display the build number and version (if any)
var versionc = &cobra.Command{
	Use:   "version",
	Short: "Show build and version",
	Run: func(c *cobra.Command, args []string) {
		fmt.Printf("Build: %s\nVersion: %s\nBuild Date: %s\n", Build, Version, BuildDate)
	},
}

// New command that will create a new empty quokka template
var newc = &cobra.Command{
	Use:   "new [output] <options>",
	Short: "Create a new quokka template",
	Args:  cobra.MinimumNArgs(1),
	Run: func(c *cobra.Command, args []string) {
		cmd.NewQuokkaTemplate(
			args[0],
			viper.GetString("name"),
			viper.GetString("description"),
			viper.GetString("version"),
			viper.GetBool("yes"),
			viper.GetBool("debug"),
		)
	},
}

func main() {
	// Initialize Cobra and Viper
	cobra.OnInitialize(cmd.Initialize)
	cmd.AddRendererFlags(rootc)
	cmd.AddGlobalFlags(rootc)

	// Add extra flags
	rootc.AddCommand(versionc)
	cmd.AddNewFlags(newc)
	rootc.AddCommand(newc)

	// Run the command
	if err := rootc.Execute(); err != nil {
		log.Fatal("Couldn't start program:", err)
	}
}
