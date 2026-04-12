package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"payment-service/internal/repository"
	"payment-service/internal/transport/http"
	"payment-service/internal/usecase"
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:Adiktop4ik@localhost:5432/payment_db?sslmode=disable")
	if err != nil {
		panic(err)
	}

	repo := &repository.PostgresRepo{DB: db}
	uc := &usecase.PaymentUseCase{Repo: repo}
	handler := &http.PaymentHandler{UC: uc}

	r := gin.Default()
	r.POST("/payments", handler.HandleAuthorize)
	r.GET("/payments/:order_id", handler.HandleGetStatus)

	r.Run(":8081")
}
