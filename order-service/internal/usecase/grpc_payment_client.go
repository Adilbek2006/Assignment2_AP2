package usecase

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

type GRPCPaymentClient struct {
	Addr string
}

func (c *GRPCPaymentClient) CreatePayment(orderID string, amount int64) (string, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)
	resp, err := client.ProcessPayment(context.Background(), &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  float64(amount),
	})

	if err != nil {
		return "Failed", err
	}

	if resp.Success {
		return "Authorized", nil
	}
	return "Declined", nil
}
