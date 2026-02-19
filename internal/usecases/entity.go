package usecases

import (
	"context"
	"github.com/dontpanicw/SalesTracker/internal/domain"
	"github.com/dontpanicw/SalesTracker/internal/port"
	"time"
)

type useCases struct {
	repo port.Repository
}

// New создает новый экземпляр use cases
func New(repo port.Repository) port.UseCases {
	return &useCases{repo: repo}
}

func (u *useCases) CreateItem(ctx context.Context, item *domain.Item) error {
	if err := item.Validate(); err != nil {
		return err
	}
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	return u.repo.Create(ctx, item)
}

func (u *useCases) GetItem(ctx context.Context, id int64) (*domain.Item, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *useCases) GetItems(ctx context.Context, from, to *time.Time) ([]*domain.Item, error) {
	return u.repo.GetAll(ctx, from, to)
}

func (u *useCases) UpdateItem(ctx context.Context, item *domain.Item) error {
	if err := item.Validate(); err != nil {
		return err
	}
	item.UpdatedAt = time.Now()
	return u.repo.Update(ctx, item)
}

func (u *useCases) DeleteItem(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func (u *useCases) GetAnalytics(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
	return u.repo.GetAnalytics(ctx, from, to)
}
