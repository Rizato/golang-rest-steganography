package worker

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	AMQP_URL string = os.Getenv("AMQP_URL")
)

func main() error {
	conn, err := amqp.Dial()
}
