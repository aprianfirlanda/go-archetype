package cmd

/*
Copyright © 2025 APRIAN FIRLANDA IMANI <aprianfirlanda@gmail.com>
*/

import (
	"context"
	"fmt"
	"go-archetype/internal/infrastructure/persistance/gorm/migrate"
	"strconv"

	"github.com/spf13/cobra"
)

var migrationsDir = "migrations"

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration commands",
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCreateCmd.Flags().String("type", "sql", "Migration type (sql|go)")
	migrateCmd.AddCommand(
		migrateCreateCmd,
		migrateUpCmd,
		migrateUpToCmd,
		migrateDownCmd,
		migrateDownToCmd,
		migrateStatusCmd,
		migrateVersionCmd,
	)
}

func initGooseMigrator() (*migrate.GooseMigrator, error) {
	sqlDB, err := dbConn.DB()
	if err != nil {
		logger.WithError(err).Error("failed to get sql.DB")
		return nil, err
	}

	return migrate.NewGooseMigrator(sqlDB, migrationsDir), nil
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create new migration file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mType, _ := cmd.Flags().GetString("type")
		gooseMigrator := migrate.NewGooseMigrator(nil, migrationsDir)
		return gooseMigrator.Create(args[0], mType)
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		gooseMigrator, err := initGooseMigrator()
		if err != nil {
			return err
		}
		return gooseMigrator.Up(context.Background())
	},
}

var migrateUpToCmd = &cobra.Command{
	Use:   "up-to [version]",
	Short: "Migrate up to a specific version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		gooseMigrator, err := initGooseMigrator()
		if err != nil {
			return err
		}
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		return gooseMigrator.UpTo(context.Background(), version)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback last migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		gooseMigrator, err := initGooseMigrator()
		if err != nil {
			return err
		}
		return gooseMigrator.Down(context.Background())
	},
}

var migrateDownToCmd = &cobra.Command{
	Use:   "down-to [version]",
	Short: "Rollback to a specific version",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		gooseMigrator, err := initGooseMigrator()
		if err != nil {
			return err
		}
		return gooseMigrator.DownTo(context.Background(), version)
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		gooseMigrator, err := initGooseMigrator()
		if err != nil {
			return err
		}
		return gooseMigrator.Status(context.Background())
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	RunE: func(cmd *cobra.Command, args []string) error {
		gooseMigrator, err := initGooseMigrator()
		if err != nil {
			return err
		}
		version, err := gooseMigrator.Version(context.Background())
		if err != nil {
			return err
		}
		fmt.Println(version)
		return nil
	},
}
