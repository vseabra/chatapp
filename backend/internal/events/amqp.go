package events

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AMQP struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewAMQP(ctx context.Context, uri string, exchange string) (*AMQP, error) {
	deadline := time.Now().Add(60 * time.Second)
	var lastErr error
	for attempt := 0; time.Now().Before(deadline); attempt++ {
		conn, err := amqp.Dial(uri)
		if err == nil {
			ch, chErr := conn.Channel()
			if chErr == nil {
				if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err == nil {
					return &AMQP{conn: conn, channel: ch, exchange: exchange}, nil
				}
				lastErr = err
				_ = ch.Close()
				_ = conn.Close()
			} else {
				lastErr = chErr
				_ = conn.Close()
			}
		} else {
			lastErr = err
		}
		sleep := time.Duration(200*(1<<attempt)) * time.Millisecond
		if sleep > 3*time.Second {
			sleep = 3 * time.Second
		}
		time.Sleep(sleep)
	}
	return nil, lastErr
}

func (a *AMQP) Close() {
	if a.channel != nil {
		_ = a.channel.Close()
	}
	if a.conn != nil {
		_ = a.conn.Close()
	}
}

func (a *AMQP) PublishJSON(ctx context.Context, routingKey string, body []byte) error {
	pub := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
	}
	return a.channel.PublishWithContext(ctx, a.exchange, routingKey, false, false, pub)
}
