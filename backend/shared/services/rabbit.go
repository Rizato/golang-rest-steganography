package services

import (
	"context"
	"encoding/json"
	"steg/shared/models"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitService struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewRabbitSenderService(channel *amqp.Channel, queue amqp.Queue) *RabbitService {
	return &RabbitService{channel, queue}
}

func (service *RabbitService) Publish(ctx context.Context, jobType models.JobType, uuid uuid.UUID) error {
	message := models.NewJobMessage(jobType, uuid)
	marshalled, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return service.channel.PublishWithContext(
		ctx,
		"",                 // exchange
		service.queue.Name, // routing_key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         marshalled,
		})
}
