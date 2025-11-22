package ordergrpc

import (
	"context"

	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/loveRyujin/gorder/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	app *app.Application
}

func New(app *app.Application) *Server {
	return &Server{app: app}
}

func (s *Server) CreateOrder(context.Context, *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
func (s *Server) GetOrder(context.Context, *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	panic("implement me")
}
func (s *Server) UpdateOrder(context.Context, *orderpb.Order) (*emptypb.Empty, error) {
	panic("implement me")
}
