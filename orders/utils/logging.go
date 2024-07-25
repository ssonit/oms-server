package utils

import (
	"context"
	"time"

	pb "github.com/ssonit/common/protos/order"

	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	next OrdersService
}

func NewLoggingMiddleware(next OrdersService) OrdersService {
	return &LoggingMiddleware{next}
}

func (s *LoggingMiddleware) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {

	start := time.Now()
	defer func() {
		zap.L().Info("GetOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.GetOrder(ctx, p)

}

func (s *LoggingMiddleware) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CreateOrder", zap.Duration("took", time.Since(start)))
	}()

	return s.next.CreateOrder(ctx, p)
}
