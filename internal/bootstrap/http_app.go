package bootstrap

import (
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"

	"github.com/sirupsen/logrus"
)

type HttpApp struct {
	Config      *config.Config
	Log         *logrus.Entry
	DBPinger    output.DBPinger
	TaskService input.TaskService
}
