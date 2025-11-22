package stockgrpc

import (
	"context"

	"github.com/loveRyujin/gorder/common/genproto/stockpb"
	"github.com/loveRyujin/gorder/stock/app"
)

type Server struct {
	app *app.Application
}

func New(app *app.Application) *Server {
	return &Server{app: app}
}

func (G *Server) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G *Server) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO implement me
	panic("implement me")
}
