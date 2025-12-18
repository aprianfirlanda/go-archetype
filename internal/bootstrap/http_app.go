package bootstrap

import (
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/ports"

	"github.com/sirupsen/logrus"
)

type HttpApp struct {
	Config   *config.Config
	Log      *logrus.Entry
	DBPinger ports.DBPinger
}
