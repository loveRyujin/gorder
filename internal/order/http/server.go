package orderhttp

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/gorder/order/app"
)

type Server struct {
	app *app.Application
}

func New(app *app.Application) *Server {
	return &Server{app: app}
}

func (s *Server) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	fmt.Println(customerID)
}

func (s *Server) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	fmt.Println(customerID, orderID)
}
