package orders_gateway

import (
	"context"

	pb "github.com/ssonit/common/protos/order"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.GetOrderResponse, error)
}
