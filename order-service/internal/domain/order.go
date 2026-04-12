package domain

import "time"

type Order struct {
	ID         string
	CustomerID string
	ItemName   string
	Amount     int64
	Status     string
	CreatedAt  time.Time
}

type OrderStats struct {
	Total     int64
	Pending   int64
	Paid      int64
	Failed    int64
	Cancelled int64
}

type OrderRepository interface {
	Create(o *Order) error
	UpdateStatus(id string, status string) error
	GetByID(id string) (*Order, error)
	GetStats() (*OrderStats, error)
}
