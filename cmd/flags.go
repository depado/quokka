package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddRendererFlags adds a set of flags to make the use of the program more
// flexible
func AddRendererFlags(c *cobra.Command) {
	// c.PersistentFlags().BoolP("commands", "c", false, "execute the after commands (make sure you know what it does)")

	// General options
	c.Flags().StringP("input", "i", "", "specify an input values file to automate template rendering")
	c.Flags().BoolP("keep", "k", false, "do not delete the template when operation is complete")
	c.Flags().StringP("path", "p", "", "specify if the template is actually stored in a sub-directory of the downloaded file")
	c.Flags().StringP("output", "o", "", "specify the directory where the template should be downloaded or cloned")
	// Git options
	c.Flags().Int("git.depth", 1, "depth of git clone in case of git provider")

	// TODO: Handle auth for HTTP Provider
	// c.PersistentFlags().String("user", "", "user for auth if needed")
	// c.PersistentFlags().String("password", "", "password for auth if needed")
	if err := viper.BindPFlags(c.Flags()); err != nil {
		log.Fatal("Could not bind flags")
	}
}

// AddGlobalFlags adds the persistent flags that will be added to all the
// commands
func AddGlobalFlags(c *cobra.Command) {
	c.PersistentFlags().BoolP("yes", "y", false, "Automatically accept")
	c.PersistentFlags().Bool("debug", false, "Enable or disable debug mode")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		log.Fatal("Could not bind flags")
	}
}

// AddNewFlags will apply the flags to the
func AddNewFlags(c *cobra.Command) {
	c.Flags().StringP("name", "n", "", "name of the new template")
	c.Flags().StringP("description", "d", "", "description of the new template")
	c.Flags().StringP("version", "v", "", "version of the new template")
	if err := viper.BindPFlags(c.Flags()); err != nil {
		log.Fatal("Could not bind flags")
	}
}

// Initialize will be run when cobra finishes its initialization
func Initialize() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	viper.SetEnvPrefix("quokka")

	// Environment variables
	viper.AutomaticEnv()
}
