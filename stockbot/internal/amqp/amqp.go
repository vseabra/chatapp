package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// Client wraps a RabbitMQ connection and channel.
type Client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

// New connects to RabbitMQ and returns a client.
func New(uri string) (*Client, error) {
	deadline := time.Now().Add(60 * time.Second)
	var lastErr error
	for attempt := 0; time.Now().Before(deadline); attempt++ {
		conn, err := amqp091.Dial(uri)
		if err == nil {
			ch, chErr := conn.Channel()
			if chErr == nil {
				return &Client{conn: conn, channel: ch}, nil
			}
			lastErr = fmt.Errorf("open channel: %w", chErr)
			_ = conn.Close()
		} else {
			lastErr = fmt.Errorf("dial rabbitmq: %w", err)
		}
		sleep := time.Duration(200*(1<<attempt)) * time.Millisecond
		if sleep > 3*time.Second {
			sleep = 3 * time.Second
		}
		time.Sleep(sleep)
	}
	return nil, lastErr
}

// Close closes the channel and connection.
func (c *Client) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

// Declare sets up exchange and queue binding.
func (c *Client) Declare(exchangeName, exchangeType, queueName, routingKey string) error {
	if err := c.channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange: %w", err)
	}

	q, err := c.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare queue: %w", err)
	}

	if err := c.channel.QueueBind(q.Name, routingKey, exchangeName, false, nil); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	return nil
}

// Consume starts consuming deliveries from the queue.
func (c *Client) Consume(queueName string) (<-chan amqp091.Delivery, error) {
	msgs, err := c.channel.Consume(
		queueName,
		"",    // consumer tag
		true,  // autoAck (simple bot)
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return msgs, nil
}

// PublishJSON marshals and publishes a JSON payload.
func (c *Client) PublishJSON(exchange, routingKey string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
