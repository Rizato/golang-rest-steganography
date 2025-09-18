package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"steg/shared/models"
	"syscall"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	InvalidMessage = errors.New("invalid message")
)

func New(amqpUrl string, taskQueue string) (*RabbitClient, error) {
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		taskQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, err
	}

	return NewRabbitClient(conn, ch, q), nil
}

type RabbitClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
}

func NewRabbitClient(connection *amqp.Connection, channel *amqp.Channel, queue amqp.Queue) *RabbitClient {
	return &RabbitClient{connection, channel, queue}
}

func (pub *RabbitClient) Close() error {
	err := pub.channel.Close()
	if err != nil {
		return err
	}
	return pub.connection.Close()
}

type RabbitPublisher struct {
	client *RabbitClient
}

func NewRabbitPublisher(client *RabbitClient) *RabbitPublisher {
	return &RabbitPublisher{client}
}

func (pub *RabbitPublisher) Publish(ctx context.Context, jobType models.JobType, uuid uuid.UUID) error {
	message := models.NewJobMessage(jobType, uuid)
	marshalled, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return pub.client.channel.PublishWithContext(
		ctx,
		"",                    // exchange
		pub.client.queue.Name, // routing_key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         marshalled,
		})
}

func (pub *RabbitPublisher) Close() error {
	return pub.client.Close()
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, d amqp.Delivery) error
}

type RabbitConsumer struct {
	client  *RabbitClient
	handler MessageHandler
}

func NewRabbitConsumer(client *RabbitClient, handler MessageHandler) *RabbitConsumer {
	return &RabbitConsumer{client, handler}
}

func (consumer *RabbitConsumer) SetHandler(handler MessageHandler) {
	consumer.handler = handler
}

func (consumer *RabbitConsumer) Listen(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err := consumer.client.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := consumer.client.channel.Consume(
		consumer.client.queue.Name, // queue
		"",                         // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return nil
		case d := <-msgs:
			log.Printf("Received a message: %s", d.Body)
			err = consumer.ProcessMessage(ctx, d)
			// Exit on error, assume restarting will fix it
			if err != nil {
				return err
			}
		}
	}
}

func (consumer *RabbitConsumer) ProcessMessage(ctx context.Context, msg amqp.Delivery) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()
	err := consumer.handler.HandleMessage(ctx, msg)
	if err != nil && !errors.Is(err, InvalidMessage) {
		msg.Nack(false, true)
		return err
	}
	msg.Ack(false)
	return nil
}
