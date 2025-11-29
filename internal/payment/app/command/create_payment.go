package command

import (
	"context"

	"github.com/loveRyujin/gorder/common/decorator"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	domain "github.com/loveRyujin/gorder/payment/domain/payment"
	"github.com/sirupsen/logrus"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.CommandHandler[CreatePayment, string]

type createPaymentHandler struct {
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(
	processor domain.Processor,
	orderGRPC OrderService,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CreatePaymentHandler {
	return decorator.ApplyCommandDecorators[CreatePayment, string](
		createPaymentHandler{processor: processor, orderGRPC: orderGRPC},
		logger,
		metricClient,
	)
}

func (h createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	link, err := h.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for order: %s success, payment link: %s", cmd.Order.ID, link)
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}
	if err = h.orderGRPC.UpdateOrder(ctx, newOrder); err != nil {
		return "", err
	}

	return link, err
}
