package main

import (
	"database/sql"
	"net/http"
	"order-service/internal/repository"
	orderHttp "order-service/internal/transport/http"
	"order-service/internal/usecase"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	db, _ := sql.Open("postgres", "postgres://postgres:Adiktop4ik@localhost:5432/order_db?sslmode=disable")

	httpClient := &http.Client{Timeout: 2 * time.Second}

	paymentURL := os.Getenv("PAYMENT_URL")
	if paymentURL == "" {
		paymentURL = "http://localhost:8081"
	}

	repo := &repository.PostgresRepo{DB: db}
	paymentClient := &usecase.HTTPPaymentClient{
		Client: httpClient,
		URL:    paymentURL,
	}

	uc := &usecase.OrderUseCase{
		Repo:          repo,
		PaymentClient: paymentClient,
	}

	handler := &orderHttp.OrderHandler{UC: uc}

	r := gin.Default()
	r.POST("/orders", handler.Create)
	r.GET("/orders/:id", handler.Get)
	r.PATCH("/orders/:id/cancel", handler.Cancel)
	r.GET("/orders/stats", handler.GetStats)

	r.Run(":8080")
}
