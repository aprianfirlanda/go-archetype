package healthsvc

import (
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"
)

type service struct {
	db portout.DBPinger
	// future:
	// cache portout.CachePinger
	// mq    portout.MQPinger
}

func New(db portout.DBPinger) portin.HealthService {
	return &service{
		db: db,
	}
}
