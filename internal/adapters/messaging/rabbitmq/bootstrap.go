package messagingrmq

import (
	"go-archetype/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
)

type RabbitMQ struct {
	Publisher *Publisher
	Consumer  *Consumer
}

func NewRabbitMQ(cfg config.RabbitMQ, logger *logrus.Entry) (*RabbitMQ, error) {
	conn, err := NewConnection(cfg.URL)
	if err != nil {
		return nil, err
	}

	pub, err := NewPublisher(conn)
	if err != nil {
		return nil, err
	}

	con, err := NewConsumer(conn, logger, cfg.Consumer.Retry)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Publisher: pub,
		Consumer:  con,
	}, nil
}
