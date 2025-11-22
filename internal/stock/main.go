package main

import (
	"github.com/loveRyujin/gorder/common/config"
	"github.com/loveRyujin/gorder/common/genproto/stockpb"
	grpcserver "github.com/loveRyujin/gorder/common/server/grpc"
	stockgrpc "github.com/loveRyujin/gorder/stock/grpc"
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
	grpcserver.Run(serviceName, func(server *grpc.Server) {
		service := &stockgrpc.Server{}
		stockpb.RegisterStockServiceServer(server, service)
	})
}
