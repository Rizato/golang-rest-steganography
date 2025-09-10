package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func New(amqpUrl string, taskQueue string) (*amqp.Connection, *amqp.Channel, *amqp.Queue, error) {
	// connect to rabbitmq
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, nil, nil, err
	}

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		return conn, nil, nil, err
	}

	// Declare a QUEUE
	q, err := ch.QueueDeclare(
		taskQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return conn, ch, &q, err
}
