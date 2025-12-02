package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/gorder/common/broker"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	domain "github.com/loveRyujin/gorder/payment/domain/payment"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

type Server struct {
	channel *amqp.Channel
}

func New(ch *amqp.Channel) *Server {
	return &Server{
		channel: ch,
	}
}

func (s *Server) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", s.handleWebhook)
}

func (s *Server) handleWebhook(c *gin.Context) {
	logrus.Info("Got webhook from stripe")

	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Infof("Error reading request body: %s\n", err.Error())
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	event, err := webhook.ConstructEventWithOptions(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("stripe-endpoint-secret"), webhook.ConstructEventOptions{IgnoreAPIVersionMismatch: true})
	if err != nil {
		logrus.Infof("Error verifying webhook signature: %s\n", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			logrus.Infof("error unmarshal event.data.raw into session, err = %s", err.Error())
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Infof("payment for checkout session %v success!", session.ID)

			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			marshalledOrder, err := json.Marshal(&domain.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: session.Metadata["paymentLink"],
				Items:       items,
			})
			if err != nil {
				logrus.Infof("error marshal domain.order, err = %s", err.Error())
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}

			_ = s.channel.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         marshalledOrder,
			})
			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
		}
	}
	c.JSON(http.StatusOK, nil)
}
