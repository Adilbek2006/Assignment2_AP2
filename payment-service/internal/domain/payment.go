package domain

type Payment struct {
	ID            string
	OrderID       string
	TransactionID string
	Amount        int64
	Status        string
}

type PaymentRepository interface {
	Save(p *Payment) error
	GetByOrderID(orderID string) (*Payment, error)
	ListByStatus(status string) ([]*Payment, error)
}
