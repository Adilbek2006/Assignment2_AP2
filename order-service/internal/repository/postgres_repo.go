package repository

import (
	"database/sql"
	"order-service/internal/domain"
)

type PostgresRepo struct {
	DB *sql.DB
}

func (r *PostgresRepo) Create(o *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.DB.Exec(query, o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status, o.CreatedAt)
	return err
}

func (r *PostgresRepo) UpdateStatus(id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, status, id)
	return err
}

func (r *PostgresRepo) GetByID(id string) (*domain.Order, error) {
	o := &domain.Order{}
	query := `SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id = $1`
	err := r.DB.QueryRow(query, id).Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt)
	return o, err
}

func (r *PostgresRepo) GetStats() (*domain.OrderStats, error) {
	stats := &domain.OrderStats{}
	query := `SELECT
    COUNT(*) AS total,
    SUM(CASE WHEN status = 'Pending' THEN 1 ELSE 0 END) AS pending,
    SUM(CASE WHEN status = 'Paid' THEN 1 ELSE 0 END) AS paid,
    SUM(CASE WHEN status = 'Failed' THEN 1 ELSE 0 END) AS failed,
    SUM(CASE WHEN status = 'Cancelled' THEN 1 ELSE 0 END) AS cancelled
FROM orders;`

	row := r.DB.QueryRow(query)
	err := row.Scan(&stats.Total, &stats.Pending, &stats.Paid, &stats.Failed, &stats.Cancelled)
	if err == sql.ErrNoRows {
		return &domain.OrderStats{}, nil
	}
	return stats, err
}
