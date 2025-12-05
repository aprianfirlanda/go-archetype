/*
Copyright © 2025 APRIAN FIRLANDA IMANI <aprianfirlanda@gmail.com>
*/
package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-archetype/internal/config"
	"go-archetype/internal/logging"
	"os"
)

var (
	appName = "go-archetype"
	cfgFile string
	cfg     *config.Config
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
			err := config.Initialize(appName, cfgFile, cmd)
			if err != nil {
				return err
			}
			cfg, err = config.Load(appName)
			if err != nil {
				return err
			}

			logger = logging.NewLogger(cfg.Log)

			logger.WithFields(logrus.Fields{
				"config_file":  viper.ConfigFileUsed(),
				"config_value": cfg,
			}).Debug("Configuration loaded")

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
	rootCmd.PersistentFlags().String("services-general-baseurl", "http://localhost:8080", "Base URL of the General Service API")
	rootCmd.PersistentFlags().String("services-general-apikey", "Entwt4/uwQtnD2LqMdho4cPTmFEkGAzBytS4UsO0f8g=", "SHA-256 hashed API key for General Service authentication")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
