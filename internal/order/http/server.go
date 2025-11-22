package orderhttp

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (s *Server) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	fmt.Println(customerID)
}

func (s *Server) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	fmt.Println(customerID, orderID)
}
