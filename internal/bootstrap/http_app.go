package bootstrap

import (
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/ports/outbound"

	"github.com/sirupsen/logrus"
)

type HttpApp struct {
	Config   *config.Config
	Log      *logrus.Entry
	DBPinger outbound.DBPinger
}
