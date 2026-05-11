package main

import (
	"context"
	"log"
	"time"

	pb "github.com/Adilbek2006/grpc-generated/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	log.Println("Sending ProcessPayment request to Payment Service")

	res, err := client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId:       "ORDER-556",
		Amount:        5000.0,
		PaymentMethod: "card",
	})
	if err != nil {
		log.Fatalf("Error calling ProcessPayment: %v", err)
	}

	log.Printf("Response received! Success: %v, Message: %s", res.Success, res.Message)
}
