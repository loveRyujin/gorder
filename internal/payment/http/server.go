package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func (s *Server) RegisterRoutes(c *gin.Engine) {
	c.POST("/api/webhook", s.handleWebhook)
}

func (s *Server) handleWebhook(c *gin.Context) {
	logrus.Info("Got webhook from stripe")
}
