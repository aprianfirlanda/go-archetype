/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	taskhandlermq "go-archetype/internal/adapters/messaging/rabbitmq/handler/task"
	tasksvc "go-archetype/internal/application/task/service"
	"go-archetype/internal/bootstrap"
	gorminfra "go-archetype/internal/infrastructure/persistance/gorm"
	taskgorm "go-archetype/internal/infrastructure/persistance/gorm/task"

	"github.com/spf13/cobra"
)

// consumerCmd represents the consumer command
var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Run asynchronous message consumers",
	Long: `Start background workers that consume messages from RabbitMQ and process them asynchronously.

This command initializes all required dependencies such as database connections, logging, and messaging infrastructure,
then subscribes to configured topics/queues and executes application handlers for each incoming message.

Consumers are typically used for:
- Processing domain events (e.g., task.created, task.updated)
- Running background jobs and workflows
- Decoupling heavy or non-blocking operations from the HTTP layer

This command runs continuously and should be deployed as a separate worker service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		// Infrastructure
		taskRepo := taskgorm.New(dbConn)
		uow := gorminfra.NewUnitOfWork(dbConn)

		// Application
		taskService := tasksvc.New(uow, taskRepo, rmq.Publisher)

		// Handlers
		taskHandler := taskhandlermq.New(taskService)

		// Registry
		registry := bootstrap.NewConsumerRegistry()
		registry.Register("task.created", taskHandler.Create)

		// Start all consumers
		if err := registry.Start(ctx, rmq.Consumer); err != nil {
			return err
		}

		// Block
		select {}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}
