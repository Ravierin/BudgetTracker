package repository

import (
	"github.com/Ravierin/BudgetTracker/backend/internal/model"
	"github.com/Ravierin/BudgetTracker/backend/pkg/database"
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type PositionRepository struct {
	db *database.Database
}

func NewPositionRepository(db *database.Database) *PositionRepository {
	return &PositionRepository{db: db}
}

func (r *PositionRepository) SavePosition(ctx context.Context, position model.Position) error {
	query := `
		INSERT INTO position (
			order_id, exchange, symbol, volume,
			leverage, closed_pnl, side, date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (order_id) DO UPDATE SET
			volume = EXCLUDED.volume,
			leverage = EXCLUDED.leverage,
			closed_pnl = EXCLUDED.closed_pnl,
			side = EXCLUDED.side,
			date = EXCLUDED.date,
			updated_at = NOW()
	`

	_, err := r.db.Pool.Exec(ctx, query,
		position.OrderID,
		position.Exchange,
		position.Symbol,
		position.Volume,
		position.Leverage,
		position.ClosedPnl,
		position.Side,
		position.UpdatedAt,
	)

	return err
}

func (r *PositionRepository) SavePositionBatch(ctx context.Context, positions []model.Position) error {
	batch := &pgx.Batch{}

	query := `
		INSERT INTO position (
			order_id, exchange, symbol, volume,
			leverage, closed_pnl, side, date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT (order_id) DO UPDATE SET
			volume = EXCLUDED.volume,
			leverage = EXCLUDED.leverage,
			closed_pnl = EXCLUDED.closed_pnl,
			side = EXCLUDED.side,
			date = EXCLUDED.date,
			updated_at = NOW()
	`

	for _, p := range positions {
		batch.Queue(query,
			p.OrderID,
			p.Exchange,
			p.Symbol,
			p.Volume,
			p.Leverage,
			p.ClosedPnl,
			p.Side,
			p.UpdatedAt,
		)
	}

	br := r.db.Pool.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}

func (r *PositionRepository) GetAllPositions(ctx context.Context) ([]model.Position, error) {
	query := `
		SELECT id, order_id, exchange, symbol, volume,
		       leverage, closed_pnl, side, date
		FROM position
		ORDER BY date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []model.Position
	for rows.Next() {
		var p model.Position
		err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.Exchange,
			&p.Symbol,
			&p.Volume,
			&p.Leverage,
			&p.ClosedPnl,
			&p.Side,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}

	return positions, rows.Err()
}

func (r *PositionRepository) GetPositionsByExchange(ctx context.Context, exchange string) ([]model.Position, error) {
	query := `
		SELECT id, order_id, exchange, symbol, volume,
		       leverage, closed_pnl, side, date
		FROM position
		WHERE exchange = $1
		ORDER BY date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, exchange)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []model.Position
	for rows.Next() {
		var p model.Position
		err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.Exchange,
			&p.Symbol,
			&p.Volume,
			&p.Leverage,
			&p.ClosedPnl,
			&p.Side,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}

	return positions, rows.Err()
}

func (r *PositionRepository) GetPositionsByDateRange(ctx context.Context, start, end time.Time) ([]model.Position, error) {
	query := `
		SELECT id, order_id, exchange, symbol, volume,
		       leverage, closed_pnl, side, date
		FROM position
		WHERE date BETWEEN $1 AND $2
		ORDER BY date DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []model.Position
	for rows.Next() {
		var p model.Position
		err := rows.Scan(
			&p.ID,
			&p.OrderID,
			&p.Exchange,
			&p.Symbol,
			&p.Volume,
			&p.Leverage,
			&p.ClosedPnl,
			&p.Side,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		positions = append(positions, p)
	}

	return positions, rows.Err()
}

func (r *PositionRepository) GetPositionByOrderID(ctx context.Context, orderID string) (*model.Position, error) {
	query := `
		SELECT id, order_id, exchange, symbol, volume,
		       leverage, closed_pnl, side, date
		FROM position
		WHERE order_id = $1
	`

	var p model.Position
	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&p.ID,
		&p.OrderID,
		&p.Exchange,
		&p.Symbol,
		&p.Volume,
		&p.Leverage,
		&p.ClosedPnl,
		&p.Side,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *PositionRepository) DeletePosition(ctx context.Context, id int) error {
	query := `DELETE FROM position WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}
