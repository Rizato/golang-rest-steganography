package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"steg/server/api/v1/views"
	"steg/server/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
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
	conn, err := amqp.Dial(AmqpUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}
	// v1 mux
	apiV1 := views.ConfigureViews(dbPool, ch, q)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/", apiV1)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Working!\n")
	})

	// add middleware
	handler := middleware.ConfigureMiddleware(mux)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
