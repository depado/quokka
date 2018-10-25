package cmd

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddRendererFlags adds a set of flags to make the use of the program more
// flexible
func AddRendererFlags(c *cobra.Command) {
	c.PersistentFlags().StringP("output", "o", "output", "output directory")
	c.PersistentFlags().StringP("path", "p", "", "specify if the template is actually stored in a sub-directory of the downloaded file")
	c.PersistentFlags().BoolP("commands", "c", false, "execute the after commands (make sure you know what it does)")
	c.PersistentFlags().BoolP("keep", "k", false, "do not delete the template when operation is complete")
	c.PersistentFlags().Int("git.depth", 1, "depth of git clone in case of git provider")
	c.PersistentFlags().String("template.output", "", "specify output directory for the template")
	// TODO: Handle auth for HTTP Provider
	// c.PersistentFlags().String("user", "", "user for auth if needed")
	// c.PersistentFlags().String("password", "", "password for auth if needed")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddLoggerFlags").Fatal("Couldn't bind flags")
	}
}

// Initialize will be run when cobra finishes its initialization
func Initialize() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	viper.SetEnvPrefix("projectmpl")
	// Environment variables
	viper.AutomaticEnv()
}
