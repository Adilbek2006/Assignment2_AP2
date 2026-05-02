package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	"payment-service/internal/domain"
	"payment-service/internal/repository"
	grpcTransport "payment-service/internal/transport/grpc"
	"payment-service/internal/usecase"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

type RabbitMQPublisher struct {
	channel *amqp.Channel
}

type PaymentEventPayload struct {
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

func (p *RabbitMQPublisher) PublishPaymentEvent(payment *domain.Payment) error {
	q, err := p.channel.QueueDeclare(
		"payment.completed", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	payload := PaymentEventPayload{
		OrderID:       payment.OrderID,
		Amount:        payment.Amount,
		CustomerEmail: "user@example.com",
		Status:        payment.Status,
	}
	body, _ := json.Marshal(payload)

	err = p.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		})

	if err == nil {
		log.Printf("Published event to RabbitMQ for Order #%s", payment.OrderID)
	}
	return err
}

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

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	publisher := &RabbitMQPublisher{channel: ch}

	repo := &repository.PostgresRepo{DB: db}
	uc := &usecase.PaymentUseCase{Repo: repo, Publisher: publisher}
	handler := &grpcTransport.PaymentHandler{UC: uc}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Network error: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(LoggingInterceptor))
	pb.RegisterPaymentServiceServer(grpcServer, handler)

	go func() {
		log.Printf("Payment gRPC server started on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Gracefully shutting down Payment Service")
	grpcServer.GracefulStop()
}
