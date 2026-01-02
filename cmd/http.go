package cmd

/*
Copyright © 2025 APRIAN FIRLANDA IMANI <aprianfirlanda@gmail.com>
*/

import (
	"go-archetype/internal/adapters/http/server"
	taskapp "go-archetype/internal/application/task/service"
	"go-archetype/internal/bootstrap"
	"go-archetype/internal/infrastructure/persistance/gorm"
	"go-archetype/internal/infrastructure/persistance/gorm/task"

	"github.com/spf13/cobra"
)

// @title           Go Archetype
// @version         1.0
// @description     REST API for Go Archetype

// ===== Security Definitions =====

// @securityDefinitions.apikey JWTAuth
// @in header
// @name Authorization
// @description JWT Authorization header. Format: Bearer {jwt}

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API Key Authorization header. Format: Bearer {api_key}

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Infrastructure
		dbPinger := gorm.NewPinger(dbConn)
		taskRepo := taskgorm.NewRepository(dbConn)
		uow := gorm.NewUnitOfWork(dbConn)

		// Application
		taskService := taskapp.NewService(uow, taskRepo)

		return server.StartServer(bootstrap.HttpApp{
			Config:      cfg,
			Log:         logger,
			DBPinger:    dbPinger,
			TaskService: taskService,
		})
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// httpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	httpCmd.Flags().Int("http-port", 8080, "Port to run the server on")
	httpCmd.Flags().String("jwt-secret", "your-jwt-secret", "Secret key used to sign and validate JWT tokens")
}
