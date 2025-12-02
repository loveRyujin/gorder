package service

import (
	"context"

	ordergrpc "github.com/loveRyujin/gorder/common/client/order"
	"github.com/loveRyujin/gorder/common/metrics"
	"github.com/loveRyujin/gorder/payment/adapters/grpc"
	"github.com/loveRyujin/gorder/payment/app"
	"github.com/loveRyujin/gorder/payment/app/command"
	domain "github.com/loveRyujin/gorder/payment/domain/payment"
	"github.com/loveRyujin/gorder/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (*app.Application, func()) {
	orderGRPCClient, closeOrderGRPCClient, err := ordergrpc.NewCli(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := grpc.NewOrderGRPC(orderGRPCClient)
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))

	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderGRPCClient()
	}
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) *app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := &metrics.TodoMetrics{}
	return &app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
