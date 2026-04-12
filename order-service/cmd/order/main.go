package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"order-service/internal/repository"
	orderGrpc "order-service/internal/transport/grpc"
	orderHttp "order-service/internal/transport/http"
	"order-service/internal/usecase"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

func startGRPCStreaming(port string, repo *repository.PostgresRepo) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Network error: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderTrackingServiceServer(grpcServer, &orderGrpc.OrderStreamHandler{Repo: repo})

	log.Printf("Order gRPC Streaming started on port %s", port)
	grpcServer.Serve(lis)
}

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		panic(err)
	}

	repo := &repository.PostgresRepo{DB: db}

	paymentClient := &usecase.GRPCPaymentClient{
		Addr: os.Getenv("PAYMENT_GRPC_ADDR"),
	}

	uc := &usecase.OrderUseCase{
		Repo:          repo,
		PaymentClient: paymentClient,
	}

	go startGRPCStreaming(os.Getenv("ORDER_GRPC_PORT"), repo)

	handler := &orderHttp.OrderHandler{UC: uc}
	r := gin.Default()
	r.POST("/orders", handler.Create)
	r.GET("/orders/:id", handler.Get)
	r.PATCH("/orders/:id/cancel", handler.Cancel)
	r.GET("/orders/stats", handler.GetStats)

	httpPort := os.Getenv("HTTP_PORT")
	log.Printf("Order REST API started on port %s", httpPort)
	r.Run(":" + httpPort)
}
