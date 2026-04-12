package usecase

type PaymentClient interface {
	CreatePayment(orderID string, amount int64) (string, error)
}
