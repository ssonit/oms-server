package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/ssonit/common"
	"github.com/ssonit/common/discovery"
	"github.com/ssonit/common/discovery/consul"
	"github.com/ssonit/oms-orders/biz"
	"github.com/ssonit/oms-orders/storage"
	"github.com/ssonit/oms-orders/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpcHandler "github.com/ssonit/oms-orders/transport/grpc"
)

var (
	serviceName      = "orders"
	orderServiceAddr = common.EnvConfig("ORDER_SERVICE_ADDR", "localhost:50051")
	consulAddr       = common.EnvConfig("CONSUL_ADDR", "localhost:8500")
	mongoUser        = common.EnvConfig("MONGO_DB_USERNAME", "root")
	mongoPass        = common.EnvConfig("MONGO_DB_PASSWORD", "admin")
	mongoAddr        = common.EnvConfig("MONGO_DB_HOST", "localhost:27017")
)

func connectMongoDB(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())

	return client, err
}

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, orderServiceAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logger.Error("Failed to health check", zap.Error(err))
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", orderServiceAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	defer lis.Close()

	uri := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, mongoPass, mongoAddr)
	fmt.Println(uri)
	mongoClient, err := connectMongoDB(uri)
	if err != nil {
		logger.Fatal("failed to connect to mongo db", zap.Error(err))
	}

	store := storage.NewStore(mongoClient)
	service := biz.NewService(store)
	serviceWithLogging := utils.NewLoggingMiddleware(service)

	grpcHandler.NewGRPCHandler(grpcServer, serviceWithLogging)

	logger.Info("Server grpc listening on ", zap.String("port", orderServiceAddr))

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
