package grpc

import (
	"context"

	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/sirupsen/logrus"
)

type OrderGRPC struct {
	client orderpb.OrderServiceClient
}

func NewOrderGRPC(client orderpb.OrderServiceClient) *OrderGRPC {
	return &OrderGRPC{client: client}
}

func (o OrderGRPC) UpdateOrder(ctx context.Context, order *orderpb.Order) error {
	_, err := o.client.UpdateOrder(ctx, order)
	if err != nil {
		logrus.Warnf("payment_adapter||update_order,err=%v", err)
		return err
	}

	logrus.Info("payment_adapter||update_order successfully")
	return nil
}
