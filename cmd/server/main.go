package main

import (
	"encoding/json"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
)

func main() {
	if connection, err := amqp.Dial("amqp://guest:guest@rabbit-mq:5672/"); err != nil {
		log.Fatal("Failed to connect to RabbitMQ", "error", err)
	} else {
		conn = connection
	}
	log.Info("Connected to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel", "error", err)
	}
	requests, _ := ch.QueueDeclare(
		"requests", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	responses, _ := ch.QueueDeclare(
		"responses", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	go func() {
		messages, _ := ch.Consume(
			responses.Name, // queue
			"",             // consumer
			true,           // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		for msg := range messages {
			log.Info("Received a message", "message", string(msg.Body), "requestId", msg.CorrelationId)
			msg.Ack(false)
			log.Info("Acknowledged message")
		}
	}()

	for {

		obj := map[string]any{
			"command": "v1/timezone/get",
			"nodeId":  "00000000-0000-0000-0000-000000000000",
		}
		body, _ := json.Marshal(obj)

		ch.Publish(
			requests.Name, // exchange
			"",            // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType:   "text/plain",
				Body:          body,
				CorrelationId: uuid.New().String(),
			},
		)
		time.Sleep(5 * time.Second)
	}
}
