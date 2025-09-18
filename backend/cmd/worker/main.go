package main

import (
	"context"
	"log"
	"os"
	"steg/shared/rabbit"
	"steg/shared/repositories"
	"steg/shared/services"
	"steg/worker"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	AmqpUrl     = os.Getenv("AMQP_URL")
	TaskQueue   = os.Getenv("TASK_QUEUE")
	DatabaseURL = os.Getenv("DATABASE_URL")
)

func main() {
	rabbitClient, err := rabbit.New(AmqpUrl, TaskQueue)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitClient.Close()

	// Create Handler
	dbPool, err := pgxpool.New(context.Background(), DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	handler := makeJobHandler(dbPool)
	rabbitConsumer := rabbit.NewRabbitConsumer(rabbitClient, handler)

	// Listen and respond to message
	err = rabbitConsumer.Listen(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

// TODO Use factory pattern?
func makeJobHandler(dbPool *pgxpool.Pool) *worker.JobMessageHandler {
	imageRepository := repositories.NewImageRepository(dbPool)
	embedJobRepository := repositories.NewEmbedJobRepository(dbPool)
	extractJobRepository := repositories.NewExtractJobRepository(dbPool)
	embedJobService := services.NewEmbedJobService(embedJobRepository, imageRepository)
	extractJobService := services.NewExtractJobService(extractJobRepository, imageRepository)
	return worker.NewJobMessageHandler(embedJobService, extractJobService)
}
