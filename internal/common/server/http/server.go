package http

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Run(serviceName string, wrapper func(*gin.Engine)) {
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		addr = viper.GetString("fallback-http-addr")
	}
	runHTTPServerOnAddr(addr, wrapper)
}

func runHTTPServerOnAddr(addr string, wrapper func(*gin.Engine)) {
	apiRouter := gin.New()
	wrapper(apiRouter)
	apiRouter.Group("/api")
	if err := apiRouter.Run(addr); err != nil {
		log.Printf("failed to start http server on %s\n", addr)
		panic(err)
	}
}
