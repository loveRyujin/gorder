package main

import (
	"context"

	"github.com/loveRyujin/gorder/common/config"
	"github.com/loveRyujin/gorder/common/genproto/stockpb"
	grpcserver "github.com/loveRyujin/gorder/common/server/grpc"
	stockgrpc "github.com/loveRyujin/gorder/stock/grpc"
	"github.com/loveRyujin/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Panic(err)
	}
}

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)

	switch serverType {
	case "grpc":
		grpcserver.Run(serviceName, func(server *grpc.Server) {
			svc := stockgrpc.New(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		// todo...
	default:
		panic("unsupported server type")
	}

}
