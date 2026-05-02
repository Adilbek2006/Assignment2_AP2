package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	targetStatus := "Authorized"
	log.Printf("Requesting payments with status: %s", targetStatus)

	resp, err := client.ListPayments(ctx, &pb.ListPaymentsRequest{
		Status: targetStatus,
	})
	if err != nil {
		log.Fatalf("Failed to call ListPayments: %v", err)
	}

	log.Printf("Payments found: %d", len(resp.Payments))
	for i, payment := range resp.Payments {
		statusText := "Declined"
		if payment.Success {
			statusText = "Accepted"
		}
		log.Printf("[%d] %s | %s", i+1, statusText, payment.Message)
	}
}
