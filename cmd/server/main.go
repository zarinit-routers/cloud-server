package main

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, "error", err)
	}
}

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
	requests, err := ch.QueueDeclare(
		"requests", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")
	responses, err := ch.QueueDeclare(
		"responses", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		messages, err := ch.Consume(
			responses.Name, // queue
			"",             // consumer
			true,           // auto-ack
			false,          // exclusive
			false,          // no-local
			false,          // no-wait
			nil,            // args
		)
		failOnError(err, "Failed to consume messages")
		for msg := range messages {
			log.Info("Received a message", "message", string(msg.Body), "requestId", msg.CorrelationId)
			msg.Ack(false)
			log.Info("Acknowledged message")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			obj := map[string]any{
				"command": "v1/timezone/get",
				"nodeId":  "00000000-0000-0000-0000-000000000000",
			}
			body, _ := json.Marshal(obj)

			log.Info("Sending a message", "message", string(body))
			err := ch.Publish(
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
			failOnError(err, "Failed publishing a message")
			time.Sleep(5 * time.Second)
		}
	}()
	wg.Wait()
}
