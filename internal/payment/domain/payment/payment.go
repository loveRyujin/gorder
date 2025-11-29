package payment

import (
	"context"

	"github.com/loveRyujin/gorder/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(context.Context, *orderpb.Order) (string, error)
}
