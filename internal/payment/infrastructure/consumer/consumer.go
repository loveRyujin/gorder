package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/loveRyujin/gorder/common/broker"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/loveRyujin/gorder/payment/app"
	"github.com/loveRyujin/gorder/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	app app.Application
	ch  *amqp.Channel
}

func NewConsumer(app app.Application, ch *amqp.Channel) *Consumer {
	return &Consumer{
		app: app,
		ch:  ch,
	}
}

func (c *Consumer) Listen() {
	q, err := c.ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("[Payment] fail to consume: queue=%s, err=%s", q.Name, err.Error())
	}

	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q)
		}
	}()

	select {}
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Payment receive a message from %s, msg=%v", q.Name, string(msg.Body))

	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	tr := otel.Tracer("rabbitmq")
	_, span := tr.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("failed to unmarshall msg to order, err=%s", err.Error())
		_ = msg.Nack(false, false)
		return
	}
	if _, err := c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{Order: o}); err != nil {
		// TODO: retry
		logrus.Infof("failed to create order, err=%s", err.Error())
		_ = msg.Nack(false, false)
		return
	}

	span.AddEvent("payment.created")
	if err := msg.Ack(false); err != nil {
		logrus.Warnf("Failed to ack message, err=%s", err.Error())
	}

	logrus.Info("Payment consume successfully")
}
