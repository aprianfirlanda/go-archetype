package logging

import (
	"github.com/sirupsen/logrus"
	"go-archetype/internal/config"
	"strings"
)

func NewLogger(logConfig config.Log) *logrus.Logger {
	log := logrus.New()

	// =========================
	// Set formatter
	// =========================
	format := strings.ToLower(logConfig.Format)
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
	level := strings.ToLower(logConfig.Level)

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
