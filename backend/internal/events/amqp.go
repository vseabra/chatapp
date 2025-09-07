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
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	return &AMQP{conn: conn, channel: ch, exchange: exchange}, nil
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
