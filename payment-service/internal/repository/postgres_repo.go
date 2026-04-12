package repository

import (
	"database/sql"
	"payment-service/internal/domain"
)

type PostgresRepo struct {
	DB *sql.DB
}

func (r *PostgresRepo) Save(p *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)
	return err
}

func (r *PostgresRepo) GetByOrderID(orderID string) (*domain.Payment, error) {
	p := &domain.Payment{}
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE order_id = $1`
	err := r.DB.QueryRow(query, orderID).Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	return p, err
}
