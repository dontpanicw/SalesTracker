package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yourusername/analytics-service/internal/domain"
)

type mockRepository struct {
	createFunc      func(ctx context.Context, item *domain.Item) error
	getByIDFunc     func(ctx context.Context, id int64) (*domain.Item, error)
	getAllFunc      func(ctx context.Context, from, to *time.Time) ([]*domain.Item, error)
	updateFunc      func(ctx context.Context, item *domain.Item) error
	deleteFunc      func(ctx context.Context, id int64) error
	getAnalytics    func(ctx context.Context, from, to time.Time) (*domain.Analytics, error)
}

func (m *mockRepository) Create(ctx context.Context, item *domain.Item) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, item)
	}
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*domain.Item, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockRepository) GetAll(ctx context.Context, from, to *time.Time) ([]*domain.Item, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, from, to)
	}
	return nil, nil
}

func (m *mockRepository) Update(ctx context.Context, item *domain.Item) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, item)
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockRepository) GetAnalytics(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
	if m.getAnalytics != nil {
		return m.getAnalytics(ctx, from, to)
	}
	return nil, nil
}

func TestUseCases_CreateItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *domain.Item
		mock    *mockRepository
		wantErr bool
	}{
		{
			name: "successful creation",
			item: &domain.Item{
				Type:     "income",
				Amount:   1000.00,
				Category: "Salary",
				Date:     time.Now(),
			},
			mock: &mockRepository{
				createFunc: func(ctx context.Context, item *domain.Item) error {
					item.ID = 1
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "validation error",
			item: &domain.Item{
				Type:     "income",
				Amount:   -100.00,
				Category: "Salary",
				Date:     time.Now(),
			},
			mock:    &mockRepository{},
			wantErr: true,
		},
		{
			name: "repository error",
			item: &domain.Item{
				Type:     "income",
				Amount:   1000.00,
				Category: "Salary",
				Date:     time.Now(),
			},
			mock: &mockRepository{
				createFunc: func(ctx context.Context, item *domain.Item) error {
					return errors.New("database error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(tt.mock)
			err := uc.CreateItem(context.Background(), tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCases_UpdateItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *domain.Item
		mock    *mockRepository
		wantErr bool
	}{
		{
			name: "successful update",
			item: &domain.Item{
				ID:       1,
				Type:     "expense",
				Amount:   500.00,
				Category: "Food",
				Date:     time.Now(),
			},
			mock: &mockRepository{
				updateFunc: func(ctx context.Context, item *domain.Item) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "validation error",
			item: &domain.Item{
				ID:       1,
				Type:     "invalid",
				Amount:   500.00,
				Category: "Food",
				Date:     time.Now(),
			},
			mock:    &mockRepository{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(tt.mock)
			err := uc.UpdateItem(context.Background(), tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCases_DeleteItem(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		mock    *mockRepository
		wantErr bool
	}{
		{
			name: "successful deletion",
			id:   1,
			mock: &mockRepository{
				deleteFunc: func(ctx context.Context, id int64) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "item not found",
			id:   999,
			mock: &mockRepository{
				deleteFunc: func(ctx context.Context, id int64) error {
					return errors.New("item not found")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(tt.mock)
			err := uc.DeleteItem(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCases_GetAnalytics(t *testing.T) {
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now()

	tests := []struct {
		name    string
		mock    *mockRepository
		want    *domain.Analytics
		wantErr bool
	}{
		{
			name: "successful analytics",
			mock: &mockRepository{
				getAnalytics: func(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
					return &domain.Analytics{
						Sum:        10000.00,
						Avg:        1000.00,
						Count:      10,
						Median:     950.00,
						Percentile: 1500.00,
					}, nil
				},
			},
			want: &domain.Analytics{
				Sum:        10000.00,
				Avg:        1000.00,
				Count:      10,
				Median:     950.00,
				Percentile: 1500.00,
			},
			wantErr: false,
		},
		{
			name: "repository error",
			mock: &mockRepository{
				getAnalytics: func(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
					return nil, errors.New("database error")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(tt.mock)
			got, err := uc.GetAnalytics(context.Background(), from, to)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalytics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && tt.want != nil {
				if got.Sum != tt.want.Sum || got.Avg != tt.want.Avg || got.Count != tt.want.Count {
					t.Errorf("GetAnalytics() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
