package consumer

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loveRyujin/gorder/common/broker"
	"github.com/loveRyujin/gorder/order/app"
	"github.com/loveRyujin/gorder/order/app/command"
	domain "github.com/loveRyujin/gorder/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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
	q, err := c.ch.QueueDeclare(broker.EventOrderPaid, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	err = c.ch.QueueBind(q.Name, "", broker.EventOrderPaid, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("[Order] fail to consume: queue=%s, err=%s", q.Name, err.Error())
	}

	go func() {
		for msg := range msgs {
			c.handleMessage(msg, q)
		}
	}()

	select {}
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue) {
	logrus.Infof("Order receive a message from %s, msg=%v", q.Name, string(msg.Body))

	o := &domain.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Warnf("failed to unmarshall msg to order, err=%s", err.Error())
		_ = msg.Nack(false, false)
		return
	}

	_, err := c.app.Commands.UpdateOrder.Handle(context.TODO(), command.UpdateOrder{
		Order: o,
		UpdateFn: func(ctx context.Context, o *domain.Order) (*domain.Order, error) {
			if !o.IsPaid() {
				return nil, errors.New("order is not paid")
			}
			return o, nil
		},
	})
	if err != nil {
		logrus.Warnf("Failed to update order, err=%s", err.Error())
		_ = msg.Nack(false, false)
		return
	}

	if err := msg.Ack(false); err != nil {
		logrus.Warnf("Failed to ack message, err=%s", err.Error())
	}

	logrus.Info("Order consume successfully")
}
