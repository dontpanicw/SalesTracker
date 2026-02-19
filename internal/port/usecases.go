package port

import (
	"context"
	"github.com/dontpanicw/SalesTracker/internal/domain"
	"time"
)

// UseCases определяет бизнес-логику приложения
type UseCases interface {
	CreateItem(ctx context.Context, item *domain.Item) error
	GetItem(ctx context.Context, id int64) (*domain.Item, error)
	GetItems(ctx context.Context, from, to *time.Time) ([]*domain.Item, error)
	UpdateItem(ctx context.Context, item *domain.Item) error
	DeleteItem(ctx context.Context, id int64) error
	GetAnalytics(ctx context.Context, from, to time.Time) (*domain.Analytics, error)
}
