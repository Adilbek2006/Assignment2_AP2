package grpc

import (
	"order-service/internal/domain"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/Adilbek2006/grpc-generated/proto"
)

type OrderStreamHandler struct {
	pb.UnimplementedOrderTrackingServiceServer
	Repo domain.OrderRepository
}

func (h *OrderStreamHandler) SubscribeToOrderUpdates(req *pb.OrderRequest, stream pb.OrderTrackingService_SubscribeToOrderUpdatesServer) error {
	orderID := req.GetOrderId()
	var currentStatus string

	for {
		order, err := h.Repo.GetByID(orderID)
		if err != nil {
			return status.Errorf(codes.NotFound, "Order not found")
		}

		if order.Status != currentStatus {
			currentStatus = order.Status

			update := &pb.OrderStatusUpdate{
				OrderId:   orderID,
				Status:    currentStatus,
				UpdatedAt: timestamppb.Now(),
			}

			if err := stream.Send(update); err != nil {
				return err // Клиент отключился
			}
		}

		time.Sleep(2 * time.Second)
	}
}
