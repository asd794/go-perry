package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"otel/proto"
	"otel/telemetry"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	orderClient     proto.OrderServiceClient
	inventoryClient proto.InventoryServiceClient
}

type createOrderRequest struct {
	UserID   string `json:"user_id"`
	SKU      string `json:"sku"`
	Quantity int32  `json:"quantity"`
}

type createOrderResponse struct {
	OrderID           string `json:"order_id"`
	OrderStatus       string `json:"order_status"`
	InventoryReserved bool   `json:"inventory_reserved"`
	RemainingStock    int32  `json:"remaining_stock"`
}

func (s *Server) createOrder(
	w http.ResponseWriter,
	r *http.Request,
) {

	ctx := r.Context()

	var req createOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			"invalid json",
			http.StatusBadRequest,
		)
		return
	}

	// -------------------------
	// Order Service
	// -------------------------

	orderResp, err := s.orderClient.CreateOrder(
		ctx,
		&proto.CreateOrderRequest{
			UserId:   req.UserID,
			Sku:      req.SKU,
			Quantity: req.Quantity,
		},
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	// -------------------------
	// Inventory Service
	// -------------------------

	inventoryResp, err := s.inventoryClient.ReserveStock(
		ctx,
		&proto.ReserveStockRequest{
			Sku:      req.SKU,
			Quantity: req.Quantity,
		},
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	response := createOrderResponse{
		OrderID:           orderResp.GetOrderId(),
		OrderStatus:       orderResp.GetStatus(),
		InventoryReserved: inventoryResp.GetSuccess(),
		RemainingStock:    inventoryResp.GetRemaining(),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}

func main() {
	shutdown, err := telemetry.Init("main-service")
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown(context.Background())

	// -------------------------
	// Order gRPC client
	// -------------------------

	orderConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer orderConn.Close()

	// -------------------------
	// Inventory gRPC client
	// -------------------------

	inventoryConn, err := grpc.NewClient(
		"localhost:50052",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer inventoryConn.Close()

	server := &Server{
		orderClient: proto.NewOrderServiceClient(
			orderConn,
		),
		inventoryClient: proto.NewInventoryServiceClient(
			inventoryConn,
		),
	}

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/orders",
		server.createOrder,
	)

	handler := otelhttp.NewHandler(
		mux,
		"main-http",
	)

	log.Println("Main Service listening on :8080")

	if err := http.ListenAndServe(
		":8080",
		handler,
	); err != nil {
		log.Fatal(err)
	}
}
