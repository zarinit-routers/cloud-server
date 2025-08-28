package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, "error", err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbit-mq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	log.Info("Connected to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

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
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			go func() {

				obj := map[string]any{
					"command": "v1/timezone/get",
					"nodeId":  "00000000-0000-0000-0000-000000000000",
				}
				body, _ := json.Marshal(obj)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				log.Info("Sending a message", "message", string(body))
				err := ch.PublishWithContext(ctx,
					requests.Name, // exchange
					"",            // routing key
					false,         // mandatory
					false,         // immediate
					amqp.Publishing{
						ContentType:   "application/json",
						Body:          body,
						CorrelationId: uuid.New().String(),
					},
				)
				if conn.IsClosed() {
					failOnError(fmt.Errorf("connection is closed"), "Failed publishing a message")
				}
				failOnError(err, "Failed publishing a message")
			}()
			time.Sleep(10 * time.Second)
		}
	}()
	wg.Wait()
}
