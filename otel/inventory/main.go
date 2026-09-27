package main

import (
	"context"
	"log"
	"net"
	"sync"

	"otel/proto"
	"otel/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventoryServer struct {
	proto.UnimplementedInventoryServiceServer

	mu    sync.Mutex
	stock map[string]int32
}

func (s *inventoryServer) ReserveStock(
	ctx context.Context,
	req *proto.ReserveStockRequest,
) (*proto.ReserveStockResponse, error) {

	tracer := otel.Tracer("inventory-service")

	ctx, span := tracer.Start(ctx, "inventory.reserve")
	defer span.End()

	if req.GetSku() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"sku is required",
		)
	}

	if req.GetQuantity() <= 0 {
		return nil, status.Error(
			codes.InvalidArgument,
			"quantity must be greater than zero",
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	current, exists := s.stock[req.GetSku()]
	if !exists {
		return nil, status.Error(
			codes.NotFound,
			"sku not found",
		)
	}

	if current < req.GetQuantity() {
		return &proto.ReserveStockResponse{
			Success:   false,
			Remaining: current,
		}, nil
	}

	s.stock[req.GetSku()] -= req.GetQuantity()

	remaining := s.stock[req.GetSku()]

	log.Printf(
		"[Inventory] sku=%s quantity=%d remaining=%d",
		req.GetSku(),
		req.GetQuantity(),
		remaining,
	)

	return &proto.ReserveStockResponse{
		Success:   true,
		Remaining: remaining,
	}, nil
}

func main() {
	shutdown, err := telemetry.Init("inventory-service")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(context.Background())

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	proto.RegisterInventoryServiceServer(
		server,
		&inventoryServer{
			stock: map[string]int32{
				"iphone":  10,
				"macbook": 5,
			},
		},
	)

	log.Println("Inventory Service listening on :50052")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
