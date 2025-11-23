package service

import (
	"context"

	"github.com/loveRyujin/gorder/common/metrics"
	"github.com/loveRyujin/gorder/order/adapters"
	"github.com/loveRyujin/gorder/order/app"
	"github.com/loveRyujin/gorder/order/app/command"
	"github.com/loveRyujin/gorder/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) *app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := &metrics.TodoMetrics{}

	return &app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}
