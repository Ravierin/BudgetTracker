package repository

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/pkg/database"
	"context"
	"time"
)

type WithdrawalRepository struct {
	db *database.Database
}

func NewWithdrawalRepository(db *database.Database) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (r *WithdrawalRepository) SaveWithdrawal(ctx context.Context, withdrawal model.Withdrawal) error {
	query := `
		INSERT INTO withdrawal (exchange, amount, currency, date)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		withdrawal.Exchange,
		withdrawal.Amount,
		withdrawal.Currency,
		withdrawal.CreatedAt,
	)
	return err
}

func (r *WithdrawalRepository) GetAllWithdrawals(ctx context.Context) ([]model.Withdrawal, error) {
	query := `
		SELECT id, exchange, amount, currency, date
		FROM withdrawal
		ORDER BY date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		err := rows.Scan(&w.ID, &w.Exchange, &w.Amount, &w.Currency, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	return withdrawals, rows.Err()
}

func (r *WithdrawalRepository) GetWithdrawalsByExchange(ctx context.Context, exchange string) ([]model.Withdrawal, error) {
	query := `
		SELECT id, exchange, amount, currency, date
		FROM withdrawal
		WHERE exchange = $1
		ORDER BY date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, exchange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		err := rows.Scan(&w.ID, &w.Exchange, &w.Amount, &w.Currency, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	return withdrawals, rows.Err()
}

func (r *WithdrawalRepository) GetWithdrawalsByDateRange(ctx context.Context, start, end time.Time) ([]model.Withdrawal, error) {
	query := `
		SELECT id, exchange, amount, currency, date
		FROM withdrawal
		WHERE date BETWEEN $1 AND $2
		ORDER BY date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []model.Withdrawal
	for rows.Next() {
		var w model.Withdrawal
		err := rows.Scan(&w.ID, &w.Exchange, &w.Amount, &w.Currency, &w.CreatedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	return withdrawals, rows.Err()
}

func (r *WithdrawalRepository) DeleteWithdrawal(ctx context.Context, id int) error {
	query := `DELETE FROM withdrawal WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
