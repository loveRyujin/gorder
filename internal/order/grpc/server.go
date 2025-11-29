package ordergrpc

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/loveRyujin/gorder/order/app"
	"github.com/loveRyujin/gorder/order/app/command"
	"github.com/loveRyujin/gorder/order/app/query"
	domain "github.com/loveRyujin/gorder/order/domain/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	app *app.Application
}

func New(app *app.Application) *Server {
	return &Server{app: app}
}

func (s *Server) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := s.app.Commands.CreateOrder.Handle(ctx, command.CreateOrder{
		CustomerID: request.CustomerID,
		Items:      request.Items,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (s *Server) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	order, err := s.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: request.CustomerID,
		OrderID:    request.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return order.ToProto(), nil
}

func (s *Server) UpdateOrder(ctx context.Context, request *orderpb.Order) (*emptypb.Empty, error) {
	order, err := domain.NewOrder(request.ID, request.CustomerID, request.Status, request.PaymentLink, request.Items)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, err = s.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, o *domain.Order) (*domain.Order, error) {
			return o, nil
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
