package messagingrmq

type RabbitMQ struct {
	Publisher *Publisher
	Consumer  *Consumer
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := NewConnection(url)
	if err != nil {
		return nil, err
	}

	pub, err := NewPublisher(conn)
	if err != nil {
		return nil, err
	}

	con, err := NewConsumer(conn)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		Publisher: pub,
		Consumer:  con,
	}, nil
}
