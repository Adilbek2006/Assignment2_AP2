package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"order-service/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OrderUseCase struct {
	Repo          domain.OrderRepository
	PaymentClient PaymentClient
	RedisClient   *redis.Client
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

	err = uc.Repo.UpdateStatus(o.ID, o.Status)

	uc.RedisClient.Del(context.Background(), "order:"+o.ID)

	return err
}

func (uc *OrderUseCase) CancelOrder(id string) error {
	o, err := uc.Repo.GetByID(id)
	if err != nil {
		return err
	}
	if o.Status != "Pending" {
		return errors.New("only Pending orders can be cancelled")
	}

	err = uc.Repo.UpdateStatus(id, "Cancelled")
	if err == nil {
		uc.RedisClient.Del(context.Background(), "order:"+id)
	}
	return err
}

func (uc *OrderUseCase) GetOrder(id string) (*domain.Order, error) {
	ctx := context.Background()
	cacheKey := "order:" + id

	cachedData, err := uc.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var o domain.Order
		if err := json.Unmarshal([]byte(cachedData), &o); err == nil {
			return &o, nil
		}
	}

	o, err := uc.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if jsonData, err := json.Marshal(o); err == nil {
		uc.RedisClient.Set(ctx, cacheKey, jsonData, 5*time.Minute)
	}

	return o, nil
}

func (uc *OrderUseCase) GetOrderStats() (*domain.OrderStats, error) {
	return uc.Repo.GetStats()
}
