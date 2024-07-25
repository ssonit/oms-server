package main

import (
	"context"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"

	orders_gateway "github.com/ssonit/oms-gateway/gateway/orders"

	"github.com/gin-gonic/gin"
	"github.com/ssonit/common"

	"github.com/ssonit/common/discovery"
	"github.com/ssonit/common/discovery/consul"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvConfig("HTTP_ADDR", ":3000")
	consulAddr  = common.EnvConfig("CONSUL_ADDR", "localhost:8500")
)

func main() {

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}
			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	r := gin.Default()

	ordersGateway := orders_gateway.NewGRPCGateway(registry)

	h := NewHandler(ordersGateway)
	h.RegisterRoutes(r)

	log.Printf("Starting server on %s", httpAddr)

	r.Run(httpAddr)
}
