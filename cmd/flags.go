package cmd

import (
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddLoggerFlags adds support to configure the level of the logger
func AddLoggerFlags(c *cobra.Command) {
	c.PersistentFlags().String("log.level", "info", "one of debug, info, warn, error or fatal")
	c.PersistentFlags().String("log.format", "text", "one of text or json")
	c.PersistentFlags().Bool("log.line", false, "enable filename and line in logs")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddLoggerFlags").Fatal("Couldn't bind flags")
	}
}

// AddConfigurationFlag adds support to provide a configuration file on the
// command line
func AddConfigurationFlag(c *cobra.Command) {
	c.PersistentFlags().String("conf", "", "configuration file to use")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).Fatal("Couldn't bind flags")
	}
}

// Initialize will be run when cobra finishes its initialization
func Initialize() {
	// Environment variables
	viper.AutomaticEnv()

	// Configuration file
	if viper.GetString("conf") != "" {
		viper.SetConfigFile(viper.GetString("conf"))
	} else {
		viper.SetConfigName("conf")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config/")
	}
	if err := viper.ReadInConfig(); err != nil {
		logrus.Debug("No configuration file found")
	}

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
