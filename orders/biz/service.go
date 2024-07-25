package biz

import (
	"context"

	"github.com/ssonit/oms-orders/model"
	"github.com/ssonit/oms-orders/utils"

	pb "github.com/ssonit/common/protos/order"
)

type service struct {
	store utils.OrdersStore
}

func NewService(store utils.OrdersStore) *service {
	return &service{store: store}
}

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	order, err := s.store.GetItem(ctx, map[string]interface{}{"_id": p.Id, "customerId": p.CustomerId})
	if err != nil {
		return nil, err
	}

	return &pb.GetOrderResponse{
		Id:         order.ID.Hex(),
		CustomerId: order.CustomerId,
		Items:      order.Items,
		Status:     order.Status,
	}, nil

}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	id, err := s.store.Create(ctx, model.OrdersCreation{
		CustomerId: p.CustomerId,
		Items:      p.Items,
		Status:     "pending",
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderResponse{
		Id:         id.Hex(),
		CustomerId: p.CustomerId,
		Status:     "pending",
	}, nil

}
