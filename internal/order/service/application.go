package service

import (
	"context"

	"github.com/loveRyujin/gorder/common/broker"
	stockgrpc "github.com/loveRyujin/gorder/common/client/stock"
	"github.com/loveRyujin/gorder/common/metrics"
	"github.com/loveRyujin/gorder/order/adapters"
	"github.com/loveRyujin/gorder/order/adapters/grpc"
	"github.com/loveRyujin/gorder/order/app"
	"github.com/loveRyujin/gorder/order/app/command"
	"github.com/loveRyujin/gorder/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (*app.Application, func()) {
	stockGRPCClient, closeStockGRPCClient, err := stockgrpc.NewCli(ctx)
	if err != nil {
		panic(err)
	}
	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	stockGRPC := grpc.NewStockGRPC(stockGRPCClient)

	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockGRPCClient()
		_ = ch.Close()
		_ = closeCh()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) *app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := &metrics.TodoMetrics{}

	return &app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient, ch),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
