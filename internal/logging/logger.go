package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

func NewLogger() *logrus.Logger {
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
