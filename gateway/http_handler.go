package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ssonit/common"
	orders_gateway "github.com/ssonit/oms-gateway/gateway/orders"

	pb "github.com/ssonit/common/protos/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	ordersGateway orders_gateway.OrdersGateway
}

func NewHandler(ordersGateway orders_gateway.OrdersGateway) *handler {
	return &handler{ordersGateway: ordersGateway}
}

func (h *handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)

	api := r.Group("/api")
	api.POST("/customers/:customerId/orders", h.HandleCreateOrder)
	api.GET("/customers/:customerId/orders/:orderId", h.HandleGetOrders)
}

func (h *handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (h *handler) HandleGetOrders(c *gin.Context) {
	customerId := c.Param("customerId")
	orderId := c.Param("orderId")

	orders, err := h.ordersGateway.GetOrder(c.Request.Context(), &pb.GetOrderRequest{
		Id:         orderId,
		CustomerId: customerId,
	})

	rStatus := status.Convert(err)

	if rStatus != nil {

		if rStatus.Code() != codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": rStatus.Message(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.SimpleSuccessResponse(orders))

}

func (h *handler) HandleCreateOrder(c *gin.Context) {

	customerId := c.Param("customerId")

	var items []*pb.ItemsWithQuantity

	if err := c.BindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := h.ordersGateway.CreateOrder(c.Request.Context(), &pb.CreateOrderRequest{
		CustomerId: customerId,
		Items:      items,
	})

	rStatus := status.Convert(err)

	if rStatus != nil {

		if rStatus.Code() != codes.InvalidArgument {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": rStatus.Message(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.SimpleSuccessResponse(order))
}
