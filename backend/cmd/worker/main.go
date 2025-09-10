package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"steg/shared/models"
	"steg/shared/repositories"
	"steg/shared/services"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	AmqpUrl     = os.Getenv("AMQP_URL")
	TaskQueue   = os.Getenv("TASK_QUEUE")
	DatabaseURL = os.Getenv("DATABASE_URL")
	InvalidJob  = errors.New("Invalid job")
)

func main() {
	conn, err := amqp.Dial(AmqpUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer ch.Close()

	// Declare a QUEUE
	q, err := ch.QueueDeclare(
		TaskQueue, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	// Create Handler
	dbPool, err := pgxpool.New(context.Background(), DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	imageRepository := repositories.NewImageRepository(dbPool)
	embedJobRepository := repositories.NewEmbedJobRepository(dbPool)
	extractJobRepository := repositories.NewExtractJobRepository(dbPool)
	embedJobService := services.NewEmbedJobService(embedJobRepository, imageRepository)
	extractJobService := services.NewExtractJobService(extractJobRepository, imageRepository)
	messageHandler := NewMessageHandler(embedJobService, extractJobService)

	var forever chan struct{}

	go func() {
		// Connect to DB
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			err := messageHandler.HandleMessage(d)
			if err != nil {
				if !errors.Is(err, InvalidJob) {
					// Exit on error, assume restarting will fix it
					log.Fatalf("Failed to handle a message: %s", err)
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

type MessageHandler struct {
	embedService   *services.EmbedJobService
	extractService *services.ExtractJobService
}

func NewMessageHandler(embedService *services.EmbedJobService, extractService *services.ExtractJobService) *MessageHandler {
	return &MessageHandler{embedService, extractService}
}

func (h *MessageHandler) HandleMessage(d amqp.Delivery) error {
	err := d.Ack(false)
	if err != nil {
		return err
	}
	var message models.JobMessage
	err = json.Unmarshal(d.Body, &message)
	if err != nil {
		return err
	}
	log.Printf("Received a message: %s", d.Body)
	return h.HandleJob(message)
}

func (h *MessageHandler) HandleJob(message models.JobMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	service, err := h.getService(message.Job)
	if err != nil {
		return err
	}
	job, err := service.GetJob(ctx, message.Uuid)
	if err != nil {
		return err
	}
	return service.Execute(ctx, job)
}

func (h *MessageHandler) getService(job models.JobType) (services.JobService, error) {
	switch job {
	case models.Embed:
		return h.embedService, nil
	case models.Extract:
		return h.extractService, nil
	default:
		return nil, InvalidJob
	}
}
