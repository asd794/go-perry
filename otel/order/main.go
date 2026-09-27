package main

import (
	"context"
	"log"
	"net"

	"otel/proto"
	"otel/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type orderServer struct {
	proto.UnimplementedOrderServiceServer
}

func (s *orderServer) CreateOrder(
	ctx context.Context,
	req *proto.CreateOrderRequest,
) (*proto.CreateOrderResponse, error) {

	tracer := otel.Tracer("order-service")

	ctx, span := tracer.Start(ctx, "order.create")
	defer span.End()

	if req.GetUserId() == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"user_id is required",
		)
	}

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

	log.Printf(
		"[Order] user=%s sku=%s quantity=%d",
		req.GetUserId(),
		req.GetSku(),
		req.GetQuantity(),
	)

	return &proto.CreateOrderResponse{
		OrderId: "order-001",
		Status:  "created",
	}, nil
}

func main() {
	shutdown, err := telemetry.Init("order-service")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(context.Background())

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer(
		grpc.StatsHandler(
			otelgrpc.NewServerHandler(),
		),
	)

	proto.RegisterOrderServiceServer(
		server,
		&orderServer{},
	)

	log.Println("Order Service listening on :50051")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
