package usecase

import (
	"errors"
	"github.com/google/uuid"
	"order-service/internal/domain"
	"time"
)

type OrderUseCase struct {
	Repo          domain.OrderRepository
	PaymentClient PaymentClient
}

func (uc *OrderUseCase) CreateOrder(o *domain.Order) error {
	if o.Amount <= 0 {
		return errors.New("amount must be > 0")
	}
	if o.CustomerID == "" || o.ItemName == "" {
		return errors.New("invalid input")
	}

	o.ID = uuid.New().String()
	o.Status = "Pending"
	o.CreatedAt = time.Now()

	if err := uc.Repo.Create(o); err != nil {
		return err
	}

	status, err := uc.PaymentClient.CreatePayment(o.ID, o.Amount)
	if err != nil {
		uc.Repo.UpdateStatus(o.ID, "Failed")
		return err
	}

	if status == "Authorized" {
		o.Status = "Paid"
	} else {
		o.Status = "Failed"
	}

	return uc.Repo.UpdateStatus(o.ID, o.Status)
}

func (uc *OrderUseCase) CancelOrder(id string) error {
	o, err := uc.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if o.Status != "Pending" {
		return errors.New("only Pending orders can be cancelled")
	}
	return uc.Repo.UpdateStatus(id, "Cancelled")
}

func (uc *OrderUseCase) GetOrder(id string) (*domain.Order, error) {
	return uc.Repo.GetByID(id)
}

func (uc *OrderUseCase) GetOrderStats() (*domain.OrderStats, error) {
	return uc.Repo.GetStats()
}
