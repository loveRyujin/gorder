package command

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loveRyujin/gorder/common/broker"
	"github.com/loveRyujin/gorder/common/decorator"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/loveRyujin/gorder/order/app/query"
	domain "github.com/loveRyujin/gorder/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type CreateOrder struct {
	CustomerID string
	Items      []*orderpb.ItemWithQuantity
}

type CreateOrderResult struct {
	OrderID string
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	repo      domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(
	orderRepo domain.Repository,
	stockGRPC query.StockService,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
	ch *amqp.Channel,
) CreateOrderHandler {
	if orderRepo == nil {
		panic("nil orderRepo")
	}
	if stockGRPC == nil {
		panic("nil stockGRPC")
	}
	if ch == nil {
		panic("nil channel ")
	}

	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{repo: orderRepo, stockGRPC: stockGRPC, channel: ch},
		logger,
		metricClient,
	)
}

func (h createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	validateItems, err := h.validate(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := h.repo.Create(ctx, &domain.Order{
		CustomerID: cmd.CustomerID,
		Items:      validateItems,
	})
	if err != nil {
		return nil, err
	}

	q, err := h.channel.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	err = h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
	})
	if err != nil {
		return nil, err
	}

	return &CreateOrderResult{OrderID: o.ID}, nil
}

func (h createOrderHandler) validate(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}

	items = packItems(items)
	resp, err := h.stockGRPC.CheckIfItemsInStock(ctx, items)
	if err != nil {
		return nil, err
	}

	return resp.Items, nil
}

func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}

	var res []*orderpb.ItemWithQuantity
	for id, quantity := range merged {
		res = append(res, &orderpb.ItemWithQuantity{
			ID:       id,
			Quantity: quantity,
		})
	}

	return res
}
