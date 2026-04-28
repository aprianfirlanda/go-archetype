package messagingrmq

import "github.com/rabbitmq/amqp091-go"

type Connection struct {
	Conn *amqp091.Connection
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Connection{Conn: conn}, nil
}
