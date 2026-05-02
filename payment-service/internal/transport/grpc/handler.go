package grpc

import (
	"context"
	"fmt"
	"payment-service/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
func (h *PaymentHandler) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	statusReq := req.GetStatus()

	payments, err := h.UC.GetPaymentsByStatus(statusReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list payments: %v", err)
	}

	var pbPayments []*pb.PaymentResponse
	for _, p := range payments {
		pbPayments = append(pbPayments, &pb.PaymentResponse{
			TransactionId: p.TransactionID,
			Success:       p.Status == "Authorized",
			Message:       fmt.Sprintf("Amount :%d", p.Amount),
		})
	}

	return &pb.ListPaymentsResponse{
		Payments: pbPayments,
	}, nil
}
