package usecase

import (
	"errors"
	"payment-service/internal/domain"

	"github.com/google/uuid"
)

type EventPublisher interface {
	PublishPaymentEvent(payment *domain.Payment) error
}

type PaymentUseCase struct {
	Repo      domain.PaymentRepository
	Publisher EventPublisher
}

func (uc *PaymentUseCase) Process(orderID string, amount int64) (*domain.Payment, error) {
	if orderID == "" {
		return nil, errors.New("order_id required")
	}
	if amount <= 0 {
		return nil, errors.New("invalid amount")
	}

	p := &domain.Payment{
		ID:            uuid.New().String(),
		OrderID:       orderID,
		TransactionID: uuid.New().String(),
		Amount:        amount,
	}

	if amount > 100000 {
		p.Status = "Declined"
	} else {
		p.Status = "Authorized"
	}

	err := uc.Repo.Save(p)
	if err != nil {
		return nil, err
	}

	if p.Status == "Authorized" && uc.Publisher != nil {
		err = uc.Publisher.PublishPaymentEvent(p)
		if err != nil {
			println("Failed to publish event:", err.Error())
		}
	}

	return p, err
}
func (uc *PaymentUseCase) GetStatus(orderID string) (*domain.Payment, error) {
	return uc.Repo.GetByOrderID(orderID)
}
func (uc *PaymentUseCase) GetPaymentsByStatus(status string) ([]*domain.Payment, error) {
	return uc.Repo.ListByStatus(status)
}
