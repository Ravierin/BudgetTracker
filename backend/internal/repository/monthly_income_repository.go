package repository

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/pkg/database"
	"context"
	"time"
)

type MonthlyIncomeRepository struct {
	db *database.Database
}

func NewMonthlyIncomeRepository(db *database.Database) *MonthlyIncomeRepository {
	return &MonthlyIncomeRepository{db: db}
}

func (r *MonthlyIncomeRepository) SaveMonthlyIncome(ctx context.Context, income model.MonthlyIncome) error {
	query := `
		INSERT INTO monthly_income (exchange, amount, pnl, date)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		income.Exchange,
		income.Amount,
		income.PNL,
		income.CreatedAt,
	)
	return err
}

func (r *MonthlyIncomeRepository) GetAllMonthlyIncomes(ctx context.Context) ([]model.MonthlyIncome, error) {
	query := `
		SELECT id, exchange, amount, pnl, date
		FROM monthly_income
		ORDER BY date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []model.MonthlyIncome
	for rows.Next() {
		var i model.MonthlyIncome
		err := rows.Scan(&i.ID, &i.Exchange, &i.Amount, &i.PNL, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}
	return incomes, rows.Err()
}

func (r *MonthlyIncomeRepository) GetIncomesByExchange(ctx context.Context, exchange string) ([]model.MonthlyIncome, error) {
	query := `
		SELECT id, exchange, amount, pnl, date
		FROM monthly_income
		WHERE exchange = $1
		ORDER BY date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, exchange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []model.MonthlyIncome
	for rows.Next() {
		var i model.MonthlyIncome
		err := rows.Scan(&i.ID, &i.Exchange, &i.Amount, &i.PNL, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}
	return incomes, rows.Err()
}

func (r *MonthlyIncomeRepository) GetIncomesByDateRange(ctx context.Context, start, end time.Time) ([]model.MonthlyIncome, error) {
	query := `
		SELECT id, exchange, amount, pnl, date
		FROM monthly_income
		WHERE date BETWEEN $1 AND $2
		ORDER BY date DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incomes []model.MonthlyIncome
	for rows.Next() {
		var i model.MonthlyIncome
		err := rows.Scan(&i.ID, &i.Exchange, &i.Amount, &i.PNL, &i.CreatedAt)
		if err != nil {
			return nil, err
		}
		incomes = append(incomes, i)
	}
	return incomes, rows.Err()
}

func (r *MonthlyIncomeRepository) DeleteMonthlyIncome(ctx context.Context, id int) error {
	query := `DELETE FROM monthly_income WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
