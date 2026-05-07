package logging

import (
	"go-archetype/internal/infrastructure/config"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(cfg config.Log) *logrus.Entry {
	logger := logrus.New()

	// =========================
	// Set formatter
	// =========================
	switch strings.ToLower(cfg.Format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// =========================
	// Set log level
	// =========================
	switch strings.ToLower(cfg.Level) {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetOutput(os.Stdout)

	return logrus.NewEntry(logger)
}
