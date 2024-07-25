package utils

import (
	"context"

	pb "github.com/ssonit/common/protos/order"
	"github.com/ssonit/oms-orders/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.GetOrderResponse, error)
}

type OrdersStore interface {
	Create(context.Context, model.OrdersCreation) (primitive.ObjectID, error)
	GetItem(ctx context.Context, filter map[string]interface{}) (*model.OrderItem, error)
}
