package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/gorder/common/config"
	httpserver "github.com/loveRyujin/gorder/common/server/http"
	orderhttp "github.com/loveRyujin/gorder/order/http"
	"github.com/loveRyujin/gorder/order/ports"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		logrus.Panic(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	httpserver.Run(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, &orderhttp.Server{}, ports.GinServerOptions{
			BaseURL: "/api",
		})
	})
}
