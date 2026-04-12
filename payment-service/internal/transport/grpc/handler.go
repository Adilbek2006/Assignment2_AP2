package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"payment-service/internal/usecase"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	UC *usecase.PaymentUseCase
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	payment, err := h.UC.Process(req.GetOrderId(), int64(req.GetAmount()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error: %v", err)
	}

	success := payment.Status == "Authorized"

	return &pb.PaymentResponse{
		TransactionId: payment.TransactionID,
		Success:       success,
		Message:       "Payment status: " + payment.Status,
	}, nil
}
