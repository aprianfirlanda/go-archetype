package bootstrap

import (
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"

	"github.com/sirupsen/logrus"
)

type HttpApp struct {
	Config        *config.Config
	Log           *logrus.Entry
	DBPinger      portout.DBPinger
	TaskService   portin.TaskService
	HealthService portin.HealthService
}
