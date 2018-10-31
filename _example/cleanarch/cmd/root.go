package cmd

import (
	"strings"

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

// AddServerFlags adds support to configure the server
func AddServerFlags(c *cobra.Command) {
	c.PersistentFlags().String("server.host", "127.0.0.1", "host on which the server should listen")
	c.PersistentFlags().Int("server.port", 8080, "port on which the server should listen")
	c.PersistentFlags().String("server.mode", "release", "server mode can be either 'debug', 'test' or 'release'")
	c.PersistentFlags().String("server.token", "", "authorization token to use if any")
	c.PersistentFlags().String("server.allowedOrigins", "*", "allowed origins for the server")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddServerFlags").Fatal("Couldn't bind flags")
	}
}

// AddConfigurationFlag adds support to provide a configuration file on the
// command line
func AddConfigurationFlag(c *cobra.Command) {
	c.PersistentFlags().String("conf", "", "configuration file to use")
	if err := viper.BindPFlags(c.PersistentFlags()); err != nil {
		logrus.WithError(err).WithField("step", "AddConfigurationFlag").Fatal("Couldn't bind flags")
	}
}

// Initialize will be run when cobra finishes its initialization
func Initialize() {
	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	// Set log level
	lvl := viper.GetString("log.level")
	l, err := logrus.ParseLevel(lvl)
	if err != nil {
		logrus.WithFields(logrus.Fields{"level": lvl, "fallback": "info"}).Warn("Invalid log level")
	} else {
		logrus.SetLevel(l)
	}

	// Set log format
	switch viper.GetString("log.format") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	// Defines if logrus should display filenames and line where the log occured
	if viper.GetBool("log.line") {
		logrus.AddHook(filename.NewHook())
	}
}
