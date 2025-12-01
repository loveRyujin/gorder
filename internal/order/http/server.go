package orderhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/gorder/common/genproto/orderpb"
	"github.com/loveRyujin/gorder/order/app"
	"github.com/loveRyujin/gorder/order/app/command"
	"github.com/loveRyujin/gorder/order/app/query"
)

type Server struct {
	app *app.Application
}

func New(app *app.Application) *Server {
	return &Server{app: app}
}

func (s *Server) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r, err := s.app.Commands.CreateOrder.Handle(c, command.CreateOrder{
		CustomerID: customerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "success",
		"customer_id": customerID,
		"order_id":    r.OrderID,
	})
}

func (s *Server) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	o, err := s.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": o})
}
