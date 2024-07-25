package grpc_handler

import (
	"context"

	pb "github.com/ssonit/common/protos/order"
	"github.com/ssonit/oms-orders/utils"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service utils.OrdersService
}

func NewGRPCHandler(grpcServer *grpc.Server, service utils.OrdersService) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return h.service.GetOrder(ctx, p)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {

	order, err := h.service.CreateOrder(ctx, p)
	if err != nil {
		return nil, err
	}
	return order, nil

}
