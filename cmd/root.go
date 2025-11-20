/*
Copyright © 2025 APRIAN FIRLANDA IMANI <aprianfirlanda@gmail.com>
*/
package cmd

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	appName = "go-archetype"
	cfgFile string
	logger  *logrus.Logger
	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "go-archetype",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := initializeConfig(cmd)
			if err != nil {
				return err
			}
			logger = initializeLogger()
			logger.Tracef("Configuration initialized. Using config file: %s", viper.ConfigFileUsed())

			return nil
		},
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+appName+"/config.yaml)")
	rootCmd.PersistentFlags().String("log-format", "text", "log format (text, json)")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (trace, debug, info, warn, error, fatal, panic)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initializeConfig(cmd *cobra.Command) error {
	// 1. Set up Viper to use environment variables.
	viper.SetEnvPrefix(strings.ToUpper(strings.ReplaceAll(appName, "-", "")))
	// Allow for nested keys in environment variables (e.g. `GOARCHETYPE_DATABASE_HOST`)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	// 2. Handle the configuration file.
	if cfgFile != "" {
		// Use the config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for a config file in default locations.
		home, err := os.UserHomeDir()
		// Only panic if we can't get the home directory.
		cobra.CheckErr(err)

		// Search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/." + appName)
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 3. Read the configuration file.
	// If a config file is found, read it in. We use a robust error check
	// to ignore "file not found" errors, but panic on any other error.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// 4. Bind Cobra flags to Viper.
	// This is the magic that makes the flag values available through Viper.
	// It binds the full flag set of the command passed in.
	bindSet := func(fs *pflag.FlagSet) error {
		var err error

		fs.VisitAll(func(f *pflag.Flag) {
			// convert: log-format → log.format
			key := strings.ReplaceAll(f.Name, "-", ".")

			if e := viper.BindPFlag(key, f); e != nil && err == nil {
				err = e
			}
		})

		return err
	}
	err := bindSet(cmd.Flags())
	if err != nil {
		return err
	}

	return nil
}

func initializeLogger() *logrus.Logger {
	log := logrus.New()

	// =========================
	// Set formatter
	// =========================
	format := strings.ToLower(viper.GetString("log.format"))
	switch format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// =========================
	// Set log level
	// =========================
	level := strings.ToLower(viper.GetString("log.level"))

	switch level {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}
