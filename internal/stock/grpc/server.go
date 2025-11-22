package stockgrpc

import (
	"context"

	"github.com/loveRyujin/gorder/common/genproto/stockpb"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (G *Server) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G *Server) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}
