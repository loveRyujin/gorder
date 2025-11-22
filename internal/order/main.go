package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/gorder/common/config"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	grpcserver "github.com/loveRyujin/gorder/common/server/grpc"
	httpserver "github.com/loveRyujin/gorder/common/server/http"
	ordergrpc "github.com/loveRyujin/gorder/order/grpc"
	orderhttp "github.com/loveRyujin/gorder/order/http"
	"github.com/loveRyujin/gorder/order/ports"
	"github.com/loveRyujin/gorder/order/service"
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
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application := service.NewApplication(ctx)

	go grpcserver.Run(serviceName, func(server *grpc.Server) {
		svc := ordergrpc.New(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	httpserver.Run(serviceName, func(router *gin.Engine) {
		svc := orderhttp.New(application)
		ports.RegisterHandlersWithOptions(router, svc, ports.GinServerOptions{
			BaseURL: "/api",
		})
	})
}
