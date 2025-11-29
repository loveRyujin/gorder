package query

import (
	"context"

	"github.com/loveRyujin/gorder/common/decorator"
	domain "github.com/loveRyujin/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	repo domain.Repository
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}

	return decorator.ApplyQueryDecorators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{repo: orderRepo},
		logger,
		metricClient,
	)
}

func (h getCustomerOrderHandler) Handle(ctx context.Context, query GetCustomerOrder) (*domain.Order, error) {
	o, err := h.repo.Get(ctx, query.OrderID, query.CustomerID)
	if err != nil {
		return nil, err
	}

	return o, nil
}
