package service

import (
	"context"

	stockgrpc "github.com/loveRyujin/gorder/common/client/stock"
	"github.com/loveRyujin/gorder/common/metrics"
	"github.com/loveRyujin/gorder/order/adapters"
	"github.com/loveRyujin/gorder/order/adapters/grpc"
	"github.com/loveRyujin/gorder/order/app"
	"github.com/loveRyujin/gorder/order/app/command"
	"github.com/loveRyujin/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (*app.Application, func()) {
	stockGRPCClient, closeStockGRPCClient, err := stockgrpc.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockGRPCClient)

	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockGRPCClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) *app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := &metrics.TodoMetrics{}

	return &app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
