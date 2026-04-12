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
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderTrackingServiceClient(conn)
	stream, err := client.SubscribeToOrderUpdates(context.Background(), &pb.OrderRequest{
		OrderId: "38825392-e4c1-4d84-a9eb-69ece5fe70fc",
	})
	if err != nil {
		log.Fatalf("Ошибка подписки: %v", err)
	}

	log.Println("Успешно подписались. Ожидание обновлений...")

	for {
		update, err := stream.Recv()
		if err == io.EOF {
			log.Println("Поток закрыт сервером")
			break
		}
		if err != nil {
			log.Fatalf("Ошибка чтения из потока: %v", err)
		}

		log.Printf("ОБНОВЛЕНИЕ СТАТУСА -> Заказ: %s | Статус: %s | Время: %s",
			update.OrderId, update.Status, update.UpdatedAt.AsTime().Format("15:04:05"))
	}
}
