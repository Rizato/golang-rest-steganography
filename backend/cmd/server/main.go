package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"steg/server/api/v1/views"
	"steg/server/middleware"
	"steg/shared/rabbit"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DatabaseUrl = os.Getenv("DATABASE_URL")
	AmqpUrl     = os.Getenv("AMQP_URL")
	TaskQueue   = os.Getenv("TASK_QUEUE")
)

func main() {
	log.Printf("Listening on port 8080")
	// connect to db
	dbPool, err := pgxpool.New(context.Background(), DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	// connect to rabbitmq
	rabbitClient, err := rabbit.New(AmqpUrl, TaskQueue)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitClient.Close()
	rabbitPublisher := rabbit.NewRabbitPublisher(rabbitClient)

	// v1 mux
	apiV1 := views.ConfigureViews(dbPool, rabbitPublisher)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", apiV1)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Working!\n")
	})

	// add middleware
	handler := middleware.ConfigureMiddleware(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
