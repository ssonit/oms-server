package grpc_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ssonit/oms-orders/utils"
	"google.golang.org/grpc"

	kafka "github.com/segmentio/kafka-go"
	pb "github.com/ssonit/common/protos/order"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service       utils.OrdersService
	kafkaProducer *kafka.Writer
}

func NewGRPCHandler(grpcServer *grpc.Server, service utils.OrdersService, kafkaProducer *kafka.Writer) {
	handler := &grpcHandler{
		service:       service,
		kafkaProducer: kafkaProducer,
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

	jsonBody, _ := json.Marshal(order)

	msg := kafka.Message{
		Key:   []byte(order.Id),
		Value: jsonBody,
	}

	fmt.Println("Sending message to Kafka")
	err = h.kafkaProducer.WriteMessages(ctx, msg)

	if err != nil {
		return nil, err
	}

	return order, nil

}
