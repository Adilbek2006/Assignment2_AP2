package main

import (
	"context"
	"io"
	"log"

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

	client := pb.NewOrderTrackingServiceClient(conn)
	stream, err := client.SubscribeToOrderUpdates(context.Background(), &pb.OrderRequest{
		OrderId: "56b3fb65-0213-454b-8fe1-efc9f7f3179c",
	})
	if err != nil {
		log.Fatalf("Subscription error: %v", err)
	}

	log.Println("Subscribed successfully. Waiting for updates")

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			log.Println("The stream was closed by the server")
			break
		}
		if err != nil {
			log.Fatalf("Error reading from stream: %v", err)
		}

		log.Printf("STATUS UPDATE -> Order: %s | Status: %s | Time: %s",
			update.OrderId, update.Status, update.UpdatedAt.AsTime().Format("15:04:05"))
	}
}
