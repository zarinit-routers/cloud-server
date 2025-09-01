package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func getRabbitMQUrl() (string, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		return "", fmt.Errorf("environment variable RABBITMQ_URL is not set")
	}
	return url, nil
}

var (
	Connection *amqp.Connection
	Channel    *amqp.Channel
)

const (
	requestsQueue  = "requests"
	responsesQueue = "responses"
)

func Setup() error {
	url, err := getRabbitMQUrl()
	if err != nil {
		return err
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}

	Connection = conn

	log.Info("Connected to RabbitMQ")
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	Channel = ch

	_, err = ch.QueueDeclare(
		requestsQueue, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		responsesQueue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

type JSONMap = map[string]any

type Request struct {
	Command string  `json:"command"`
	NodeId  string  `json:"nodeId"`
	Args    JSONMap `json:"args"`
}

type Response struct {
	RequestError string  `json:"requestError"`
	CommandError string  `json:"commandError"`
	Data         JSONMap `json:"data"`
}

// Response.RequestError is set by other services (for example: some required argument not specified), but returning error means that there is not response
func SendRequest(r *Request) (*Response, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), getRequestTimeout())
	defer cancel()

	requestId := uuid.New().String()

	log.Info("Sending a message", "message", string(body))
	err = Channel.PublishWithContext(ctx,
		"",            // exchange
		requestsQueue, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: requestId,
		},
	)

	messages, err := Channel.Consume(
		responsesQueue, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return nil, err
	}

	for msg := range messages {
		if msg.CorrelationId == requestId {
			msg.Ack(false)

			var response Response
			err := json.Unmarshal(msg.Body, &response)
			if err != nil {
				return nil, err
			}
			return &response, nil
		}
	}

	return nil, fmt.Errorf("something went wrong")
}

func getRequestTimeout() time.Duration {
	return 20 * time.Second
}
