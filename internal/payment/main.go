package main

import (
	"github.com/loveRyujin/gorder/common/broker"
	"github.com/loveRyujin/gorder/common/config"
	"github.com/loveRyujin/gorder/common/logging"
	httpserver "github.com/loveRyujin/gorder/common/server/http"
	paymenthttp "github.com/loveRyujin/gorder/payment/http"
	"github.com/loveRyujin/gorder/payment/infrastructure/consumer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Panic(err)
	}
}

func main() {
	stripeKey := viper.GetString("stripe-key")
	logrus.Info(stripeKey)
	serviceName := viper.GetString("payment.service-name")
	serverType := viper.GetString("payment.server-to-run")

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCh()
	}()

	go consumer.NewConsumer().Listen(ch)

	switch serverType {
	case "http":
		svc := paymenthttp.New()
		httpserver.Run(serviceName, svc.RegisterRoutes)
	case "grpc":
		// todo...
	default:
		logrus.Panic("unsupported server type")
	}
}
