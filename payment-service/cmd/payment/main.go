package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"payment-service/internal/repository"
	grpcTransport "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Printf("gRPC query : %s", info.FullMethod)
	resp, err := handler(ctx, req)
	log.Printf("Method: %s is completed by %s", info.FullMethod, time.Since(start))
	return resp, err
}

func main() {
	_ = godotenv.Load()

	db, err := sql.Open("postgres", os.Getenv("DB_DSN"))
	if err != nil {
		panic(err)
	}

	repo := &repository.PostgresRepo{DB: db}
	uc := &usecase.PaymentUseCase{Repo: repo}
	handler := &grpcTransport.PaymentHandler{UC: uc}

	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Network error: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
	pb.RegisterPaymentServiceServer(grpcServer, handler)

	log.Printf("Payment gRPC server started on port %s", port)
	grpcServer.Serve(lis)
}
