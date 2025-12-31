package bootstrap

import (
	"go-archetype/internal/infrastructure/config"
	portin "go-archetype/internal/ports/input"
	portout "go-archetype/internal/ports/output"

	"github.com/sirupsen/logrus"
)

type HttpApp struct {
	Config      *config.Config
	Log         *logrus.Entry
	DBPinger    portout.DBPinger
	TaskService portin.TaskService
}
