package ordergrpc

import (
	"context"

	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
}

func New() *Server {
	return &Server{}
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
