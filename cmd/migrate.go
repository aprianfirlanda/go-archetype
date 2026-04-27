package cmd

/*
Copyright © 2025 APRIAN FIRLANDA IMANI <aprianfirlanda@gmail.com>
*/

import (
	"context"
	"fmt"
	"go-archetype/internal/infrastructure/persistance/gorm/migrate"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var migrationsDir = "migrations"

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration management",
	Long: `Manage database schema migrations.

Supports creating, applying, rolling back, and inspecting migrations.

Migration files are stored in the migrations directory.`,
	Example: `  go-archetype migrate create init_schema
  go-archetype migrate up
  go-archetype migrate down
  go-archetype migrate status`,
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

func validateVersionArg(_ *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("requires exactly 1 argument: version")
	}

	if _, err := strconv.ParseInt(args[0], 10, 64); err != nil {
		return fmt.Errorf("version must be a valid number, got: %s", args[0])
	}

	return nil
}

func withTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

var migrateCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new migration file",
	Long:  "Create a new migration file in SQL or Go format.",
	Example: `  go-archetype migrate create add_users_table
  go-archetype migrate create add_index --type go`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mType, _ := cmd.Flags().GetString("type")

		logger.WithFields(logrus.Fields{
			"name": args[0],
			"type": mType,
		}).Info("creating migration")

		gooseMigrator := migrate.NewGooseMigrator(nil, migrationsDir)
		return gooseMigrator.Create(args[0], mType)
	},
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all pending migrations",
	Long: `Apply all pending database migrations.

This will execute all migrations that have not yet been applied
to bring the database schema to the latest version.`,
	Example: "  go-archetype migrate up",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := initGooseMigrator()
		if err != nil {
			return err
		}

		ctx, cancel := withTimeout()
		defer cancel()

		logger.Info("running migrations up")
		return m.Up(ctx)
	},
}

var migrateUpToCmd = &cobra.Command{
	Use:     "up-to [version]",
	Short:   "Apply migrations up to a specific version",
	Long:    `Apply all pending migrations up to the specified version (inclusive).`,
	Example: `  go-archetype migrate up-to 20240101010101`,
	Args:    validateVersionArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := initGooseMigrator()
		if err != nil {
			return err
		}

		version, _ := strconv.ParseInt(args[0], 10, 64) // safe now

		ctx, cancel := withTimeout()
		defer cancel()

		logger.WithField("target_version", version).
			Info("running migrations up to version")

		return m.UpTo(ctx, version)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Rollback last migration",
	Long: `Rollback the most recently applied migration.

Useful for reverting a failed deployment or incorrect schema change.`,
	Example: "  go-archetype migrate down",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := initGooseMigrator()
		if err != nil {
			return err
		}

		ctx, cancel := withTimeout()
		defer cancel()

		logger.Warn("rolling back last migration")
		return m.Down(ctx)
	},
}

var migrateDownToCmd = &cobra.Command{
	Use:     "down-to [version]",
	Short:   "Rollback migrations down to a specific version",
	Long:    `Rollback applied migrations until reaching the specified version.`,
	Example: `  go-archetype migrate down-to 20240101010101`,
	Args:    validateVersionArg,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := strconv.ParseInt(args[0], 10, 64)

		m, err := initGooseMigrator()
		if err != nil {
			return err
		}

		ctx, cancel := withTimeout()
		defer cancel()

		logger.WithField("target_version", version).
			Warn("rolling back migrations to version")

		return m.DownTo(ctx, version)
	},
}

var migrateStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show migration status",
	Long: `Display the current status of all migrations.

Shows which migrations have been applied and which are still pending.`,
	Example: "  go-archetype migrate status",
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := initGooseMigrator()
		if err != nil {
			return err
		}

		ctx, cancel := withTimeout()
		defer cancel()

		return m.Status(ctx)
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current migration version",
	Long: `Display the current migration version of the database.

This represents the latest successfully applied migration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := initGooseMigrator()
		if err != nil {
			return err
		}

		ctx, cancel := withTimeout()
		defer cancel()

		version, err := m.Version(ctx)
		if err != nil {
			return err
		}

		logger.WithField("version", version).
			Info("current migration version")

		return nil
	},
}
