package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"order-service/internal/repository"
	orderGrpc "order-service/internal/transport/grpc"
	orderHttp "order-service/internal/transport/http"
	"order-service/internal/usecase"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

func rateLimiterMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		clientIP := c.ClientIP()
		key := "ratelimit:" + clientIP

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if count > 10 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests. Please wait."})
			return
		}

		c.Next()
	}
}

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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully!")

	paymentClient := &usecase.GRPCPaymentClient{
		Addr: os.Getenv("PAYMENT_GRPC_ADDR"),
	}

	uc := &usecase.OrderUseCase{
		Repo:          repo,
		PaymentClient: paymentClient,
		RedisClient:   rdb,
	}

	go startGRPCStreaming(os.Getenv("ORDER_GRPC_PORT"), repo)

	handler := &orderHttp.OrderHandler{UC: uc}
	r := gin.Default()

	r.Use(rateLimiterMiddleware(rdb))

	r.POST("/orders", handler.Create)
	r.GET("/orders/:id", handler.Get)
	r.PATCH("/orders/:id/cancel", handler.Cancel)
	r.GET("/orders/stats", handler.GetStats)

	httpPort := os.Getenv("HTTP_PORT")
	log.Printf("Order REST API started on port %s", httpPort)
	r.Run(":" + httpPort)
}
