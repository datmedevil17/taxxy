package messaging

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// ─── Event types that flow through RabbitMQ ─────────────────
// Each event is a self-describing JSON blob.

const (
	// Exchange – single fan-out exchange; services bind their own queues.
	Exchange     = "taxxy.events"
	ExchangeType = "topic"

	// Routing keys
	EventRideRequested = "ride.requested" // Trip created → notify drivers
	EventRideAccepted  = "ride.accepted"  // Driver accepted → notify rider
	EventRideStarted   = "ride.started"   // Trip in_progress
	EventRideCompleted = "ride.completed" // Trip done → trigger payment
	EventRideCancelled = "ride.cancelled" // Trip cancelled
	EventPaymentDone   = "payment.done"   // Payment processed
)

// ─── Connection ──────────────────────────────────────────────

type Connection struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	url     string
}

// Connect dials RabbitMQ with retry logic (same pattern as DB).
func Connect(url string) *Connection {
	var conn *amqp091.Connection
	var err error

	for i := 0; i < 30; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			break
		}
		log.Printf("[MQ] connection attempt %d failed: %v – retrying in 3s…", i+1, err)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("[MQ] could not connect after 30 attempts: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[MQ] channel open failed: %v", err)
	}

	// Declare the topic exchange (idempotent).
	err = ch.ExchangeDeclare(Exchange, ExchangeType, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[MQ] exchange declare failed: %v", err)
	}

	log.Println("[MQ] connected and exchange declared")
	return &Connection{conn: conn, channel: ch, url: url}
}

// ─── Publish ─────────────────────────────────────────────────

// Publish marshals payload → JSON and sends to the topic exchange.
func (c *Connection) Publish(routingKey string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.channel.PublishWithContext(
		context.Background(), // context
		Exchange,             // exchange
		routingKey,           // routing key
		false,                // mandatory
		false,                // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent, // survives broker restart
		},
	)
}

// ─── Subscribe ───────────────────────────────────────────────

// HandlerFunc is what each service registers per routing key.
type HandlerFunc func(body []byte) error

// Subscribe creates a durable queue bound to the given routing key,
// then loops consuming messages and calling handler.
// queueName should be unique per service (e.g. "driver-service.ride.requested").
func (c *Connection) Subscribe(queueName string, routingKeys []string, handler HandlerFunc) error {
	// Declare a durable queue owned by this consumer.
	_, err := c.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Bind each routing key to the queue.
	for _, key := range routingKeys {
		if err := c.channel.QueueBind(queueName, key, Exchange, false, nil); err != nil {
			return err
		}
	}

	// Prefetch 1 – process one message at a time (fair dispatch).
	if err := c.channel.Qos(1, 0, false); err != nil {
		return err
	}

	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("[MQ] handler error on %s: %v – requeueing", queueName, err)
				msg.Nack(false, true) // requeue
			} else {
				msg.Ack(false)
			}
		}
	}()

	log.Printf("[MQ] subscribed queue=%s keys=%v", queueName, routingKeys)
	return nil
}

// ─── Close ───────────────────────────────────────────────────

func (c *Connection) Close() {
	c.channel.Close()
	c.conn.Close()
}
