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
func (r *PostgresRepo) ListByStatus(status string) ([]*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE status = $1`
	rows, err := r.DB.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		p := &domain.Payment{}
		err := rows.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
		if err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
