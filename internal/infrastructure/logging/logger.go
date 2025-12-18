package logging

import (
	"go-archetype/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
)

func New(cfg config.Log) *logrus.Entry {
	return newLogrusLogger(cfg)
}
