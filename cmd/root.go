package cmd

/*
Copyright © 2025 APRIAN FIRLANDA IMANI <aprianfirlanda@gmail.com>
*/

import (
	"go-archetype/internal/config"
	"go-archetype/internal/db"
	"go-archetype/internal/logging"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	appName = "go-archetype"
	cfgFile string
	cfg     *config.Config
	logger  *logrus.Entry
	dbConn  *gorm.DB
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

			dbConn, err = db.InitPostgres(cfg.DB, logger, []any{})

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

	rootCmd.PersistentFlags().String("db-host", "localhost", "Database host")
	rootCmd.PersistentFlags().Int("db-port", 5432, "Database port")
	rootCmd.PersistentFlags().String("db-user", "app", "Database user")
	rootCmd.PersistentFlags().String("db-password", "change_me", "Database password")
	rootCmd.PersistentFlags().String("db-name", "app", "Database name")
	rootCmd.PersistentFlags().String("db-sslmode", "disable", "Database SSL mode (disable|require|verify-ca|verify-full)")
	rootCmd.PersistentFlags().String("db-timezone", "UTC", "Database time zone (e.g. UTC)")
	rootCmd.PersistentFlags().Int("db-maxopenconns", 25, "Maximum number of open connections to the database")
	rootCmd.PersistentFlags().Int("db-maxidleconns", 25, "Maximum number of idle connections in the pool")
	rootCmd.PersistentFlags().Duration("db-connmaxlifetime", time.Hour, "Maximum amount of time a connection may be reused")
	rootCmd.PersistentFlags().Duration("db-connmaxidletime", 15*time.Minute, "Maximum amount of time a connection may be idle")
	rootCmd.PersistentFlags().String("db-loglevel", "", "GORM log level (silent|error|warn|info). Empty uses default (warn)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
