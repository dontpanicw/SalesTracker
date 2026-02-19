package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dontpanicw/SalesTracker/internal/domain"
	"github.com/dontpanicw/SalesTracker/internal/port"
	"time"

	_ "github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

// New создает новый экземпляр PostgreSQL репозитория
func New(dsn string) (port.Repository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &repository{db: db}, nil
}

func (r *repository) Create(ctx context.Context, item *domain.Item) error {
	query := `
		INSERT INTO items (type, amount, category, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRowContext(
		ctx, query,
		item.Type, item.Amount, item.Category, item.Date,
		item.CreatedAt, item.UpdatedAt,
	).Scan(&item.ID)
}

func (r *repository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	query := `
		SELECT id, type, amount, category, date, created_at, updated_at
		FROM items
		WHERE id = $1
	`
	item := &domain.Item{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID, &item.Type, &item.Amount, &item.Category,
		&item.Date, &item.CreatedAt, &item.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found")
	}
	return item, err
}

func (r *repository) GetAll(ctx context.Context, from, to *time.Time) ([]*domain.Item, error) {
	query := `
		SELECT id, type, amount, category, date, created_at, updated_at
		FROM items
		WHERE ($1::timestamp IS NULL OR date >= $1)
		  AND ($2::timestamp IS NULL OR date <= $2)
		ORDER BY date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.Item
	for rows.Next() {
		item := &domain.Item{}
		if err := rows.Scan(
			&item.ID, &item.Type, &item.Amount, &item.Category,
			&item.Date, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *repository) Update(ctx context.Context, item *domain.Item) error {
	query := `
		UPDATE items
		SET type = $1, amount = $2, category = $3, date = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.ExecContext(
		ctx, query,
		item.Type, item.Amount, item.Category, item.Date,
		item.UpdatedAt, item.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM items WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *repository) GetAnalytics(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
	query := `
		SELECT 
			COALESCE(SUM(amount), 0) as sum,
			COALESCE(AVG(amount), 0) as avg,
			COUNT(*) as count,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) as median,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) as percentile_90
		FROM items
		WHERE date >= $1 AND date <= $2
	`
	analytics := &domain.Analytics{}
	err := r.db.QueryRowContext(ctx, query, from, to).Scan(
		&analytics.Sum,
		&analytics.Avg,
		&analytics.Count,
		&analytics.Median,
		&analytics.Percentile,
	)
	return analytics, err
}

func (r *repository) Close() error {
	return r.db.Close()
}
