package port

import (
	"context"
	"time"

	"github.com/yourusername/analytics-service/internal/domain"
)

// Repository определяет интерфейс для работы с хранилищем
type Repository interface {
	Create(ctx context.Context, item *domain.Item) error
	GetByID(ctx context.Context, id int64) (*domain.Item, error)
	GetAll(ctx context.Context, from, to *time.Time) ([]*domain.Item, error)
	Update(ctx context.Context, item *domain.Item) error
	Delete(ctx context.Context, id int64) error
	GetAnalytics(ctx context.Context, from, to time.Time) (*domain.Analytics, error)
}
