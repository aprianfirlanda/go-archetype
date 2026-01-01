package demohandler

import (
	"go-archetype/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
)

type DemoHandler struct {
	log *logrus.Entry
	cfg *config.Config
}

func NewDemoHandler(log *logrus.Entry, cfg *config.Config) *DemoHandler {
	return &DemoHandler{
		log: log,
		cfg: cfg,
	}
}
