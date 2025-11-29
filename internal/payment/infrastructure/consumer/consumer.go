package consumer

import (
	"github.com/loveRyujin/gorder/common/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	ch *amqp.Channel
}

func NewConsumer(ch *amqp.Channel) *Consumer {
	return &Consumer{
		ch: ch,
	}
}

func (c *Consumer) Listen() {
	q, err := c.ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := c.ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
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
	if err := msg.Ack(false); err != nil {
		logrus.Warnf("Failed to ack message, err=%s", err.Error())
	}
}
