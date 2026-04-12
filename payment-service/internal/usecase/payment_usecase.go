package usecase

import (
	"errors"
	"github.com/google/uuid"
	"payment-service/internal/domain"
)

type PaymentUseCase struct {
	Repo domain.PaymentRepository
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
	return p, err
}

func (uc *PaymentUseCase) GetStatus(orderID string) (*domain.Payment, error) {
	return uc.Repo.GetByOrderID(orderID)
}
