package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/gorder/common/config"
	"github.com/loveRyujin/gorder/common/discovery"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/loveRyujin/gorder/common/logging"
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
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Panic(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

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
