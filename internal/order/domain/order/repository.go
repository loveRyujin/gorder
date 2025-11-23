package order

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(ctx context.Context, o *Order) (*Order, error)
	Get(ctx context.Context, id, custormerID string) (*Order, error)
	Update(
		ctx context.Context,
		o *Order,
		updateFn func(context.Context, *Order) (*Order, error),
	) error
}

type NotFoundError struct {
	OrderID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order '%s' not found", e.OrderID)
}
