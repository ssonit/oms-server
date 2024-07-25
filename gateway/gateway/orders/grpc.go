package orders_gateway

import (
	"context"
	"log"

	"github.com/ssonit/common/discovery"
	pb "github.com/ssonit/common/protos/order"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerId: p.CustomerId,
		Items:      p.Items,
	})
}

func (g *gateway) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	conn, err := discovery.ServiceConnection(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrderServiceClient(conn)

	return c.GetOrder(ctx, &pb.GetOrderRequest{
		Id:         p.Id,
		CustomerId: p.CustomerId,
	})
}
