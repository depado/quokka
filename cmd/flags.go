package cmd

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddRendererFlags adds a set of flags to make the use of the program more
// flexible
func AddRendererFlags(c *cobra.Command) {
	c.PersistentFlags().StringP("output", "o", "output", "output directory of rendered template")
	c.PersistentFlags().Int("git.depth", 1, "depth of git clone in case of git provider")
	c.PersistentFlags().Bool("template.keep", false, "do not delete the template when operation is complete")
	c.PersistentFlags().String("template.path", "", "specify output directory for the template")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddLoggerFlags").Fatal("Couldn't bind flags")
	}
}

// AddLoggerFlags adds support to configure the level of the logger
func AddLoggerFlags(c *cobra.Command) {
	c.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	c.PersistentFlags().String("log.format", "text", "one of text or json")
	c.PersistentFlags().Bool("log.line", false, "enable filename and line in logs")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddLoggerFlags").Fatal("Couldn't bind flags")
	}
}

// Initialize will be run when cobra finishes its initialization
func Initialize() {
	// Environment variables
	viper.AutomaticEnv()

	lvl := viper.GetString("log.level")
	l, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.WithField("level", lvl).Warn("Invalid log level, fallback to 'info'")
	} else {
		logrus.SetLevel(l)
	}
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	if viper.GetBool("log.line") {
		logrus.AddHook(filename.NewHook())
	}
}
