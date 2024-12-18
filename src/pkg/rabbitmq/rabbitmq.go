package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QBG-P2/Voting-System/config"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewRabbitMQ(cfg *config.Config) (*RabbitMQ, error) {
	conn, err := amqp.Dial(cfg.RabbitMQ.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare necessary exchanges and queues
	err = ch.ExchangeDeclare(
		"notifications", // name
		"topic",         // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return nil, err
	}

	// Declare the notifications queue
	_, err = ch.QueueDeclare(
		"notifications_queue", // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return nil, err
	}

	// Bind the queue to the exchange
	err = ch.QueueBind(
		"notifications_queue", // queue name
		"#",                   // routing key
		"notifications",       // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}, nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message interface{}) error {
	msg := Message{
		Type:    routingKey,
		Payload: message,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return r.channel.Publish(
		"notifications", // exchange
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (r *RabbitMQ) ConsumeMessages(ctx context.Context, queueName string, handler func(msg []byte) error) error {
	msgs, err := r.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				err = handler(msg.Body)
				if err != nil {
					log.Printf("Error handling message: %v", err)
					err = msg.Nack(false, true)
					if err != nil {
						fmt.Printf("Error nacking message: %v", err)
						return
					}

				} else {
					err = msg.Ack(false)
					if err != nil {
						fmt.Printf("Error acknowledging message: %v", err)
						return
					}
				}
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
