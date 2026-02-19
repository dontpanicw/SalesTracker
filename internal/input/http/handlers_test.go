package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/dontpanicw/SalesTracker/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type mockUseCases struct {
	createItemFunc   func(ctx context.Context, item *domain.Item) error
	getItemFunc      func(ctx context.Context, id int64) (*domain.Item, error)
	getItemsFunc     func(ctx context.Context, from, to *time.Time) ([]*domain.Item, error)
	updateItemFunc   func(ctx context.Context, item *domain.Item) error
	deleteItemFunc   func(ctx context.Context, id int64) error
	getAnalyticsFunc func(ctx context.Context, from, to time.Time) (*domain.Analytics, error)
}

func (m *mockUseCases) CreateItem(ctx context.Context, item *domain.Item) error {
	if m.createItemFunc != nil {
		return m.createItemFunc(ctx, item)
	}
	return nil
}

func (m *mockUseCases) GetItem(ctx context.Context, id int64) (*domain.Item, error) {
	if m.getItemFunc != nil {
		return m.getItemFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUseCases) GetItems(ctx context.Context, from, to *time.Time) ([]*domain.Item, error) {
	if m.getItemsFunc != nil {
		return m.getItemsFunc(ctx, from, to)
	}
	return nil, nil
}

func (m *mockUseCases) UpdateItem(ctx context.Context, item *domain.Item) error {
	if m.updateItemFunc != nil {
		return m.updateItemFunc(ctx, item)
	}
	return nil
}

func (m *mockUseCases) DeleteItem(ctx context.Context, id int64) error {
	if m.deleteItemFunc != nil {
		return m.deleteItemFunc(ctx, id)
	}
	return nil
}

func (m *mockUseCases) GetAnalytics(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
	if m.getAnalyticsFunc != nil {
		return m.getAnalyticsFunc(ctx, from, to)
	}
	return nil, nil
}

func TestHandler_CreateItem(t *testing.T) {
	tests := []struct {
		name       string
		body       interface{}
		mock       *mockUseCases
		wantStatus int
	}{
		{
			name: "successful creation",
			body: domain.Item{
				Type:     "income",
				Amount:   1000.00,
				Category: "Salary",
				Date:     time.Now(),
			},
			mock: &mockUseCases{
				createItemFunc: func(ctx context.Context, item *domain.Item) error {
					item.ID = 1
					return nil
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "invalid json",
			body:       "invalid",
			mock:       &mockUseCases{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			body: domain.Item{
				Type:     "invalid",
				Amount:   1000.00,
				Category: "Salary",
				Date:     time.Now(),
			},
			mock: &mockUseCases{
				createItemFunc: func(ctx context.Context, item *domain.Item) error {
					return errors.New("validation error")
				},
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mock)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest("POST", "/api/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateItem() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_GetItems(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mock       *mockUseCases
		wantStatus int
	}{
		{
			name:  "successful get all",
			query: "",
			mock: &mockUseCases{
				getItemsFunc: func(ctx context.Context, from, to *time.Time) ([]*domain.Item, error) {
					return []*domain.Item{
						{ID: 1, Type: "income", Amount: 1000.00, Category: "Salary", Date: time.Now()},
					}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:  "with date filters",
			query: "?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z",
			mock: &mockUseCases{
				getItemsFunc: func(ctx context.Context, from, to *time.Time) ([]*domain.Item, error) {
					return []*domain.Item{}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid date format",
			query:      "?from=invalid",
			mock:       &mockUseCases{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mock)

			req := httptest.NewRequest("GET", "/api/items"+tt.query, nil)
			w := httptest.NewRecorder()

			handler.GetItems(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetItems() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_GetItem(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mock       *mockUseCases
		wantStatus int
	}{
		{
			name: "successful get",
			id:   "1",
			mock: &mockUseCases{
				getItemFunc: func(ctx context.Context, id int64) (*domain.Item, error) {
					return &domain.Item{ID: 1, Type: "income", Amount: 1000.00}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id",
			id:         "invalid",
			mock:       &mockUseCases{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "item not found",
			id:   "999",
			mock: &mockUseCases{
				getItemFunc: func(ctx context.Context, id int64) (*domain.Item, error) {
					return nil, errors.New("not found")
				},
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mock)

			req := httptest.NewRequest("GET", "/api/items/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			w := httptest.NewRecorder()

			handler.GetItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetItem() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_DeleteItem(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		mock       *mockUseCases
		wantStatus int
	}{
		{
			name: "successful deletion",
			id:   "1",
			mock: &mockUseCases{
				deleteItemFunc: func(ctx context.Context, id int64) error {
					return nil
				},
			},
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "invalid id",
			id:         "invalid",
			mock:       &mockUseCases{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "item not found",
			id:   "999",
			mock: &mockUseCases{
				deleteItemFunc: func(ctx context.Context, id int64) error {
					return errors.New("not found")
				},
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mock)

			req := httptest.NewRequest("DELETE", "/api/items/"+tt.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.id})
			w := httptest.NewRecorder()

			handler.DeleteItem(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("DeleteItem() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestHandler_GetAnalytics(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		mock       *mockUseCases
		wantStatus int
	}{
		{
			name:  "successful analytics",
			query: "?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z",
			mock: &mockUseCases{
				getAnalyticsFunc: func(ctx context.Context, from, to time.Time) (*domain.Analytics, error) {
					return &domain.Analytics{
						Sum:        10000.00,
						Avg:        1000.00,
						Count:      10,
						Median:     950.00,
						Percentile: 1500.00,
					}, nil
				},
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing parameters",
			query:      "",
			mock:       &mockUseCases{},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid date format",
			query:      "?from=invalid&to=2024-12-31T00:00:00Z",
			mock:       &mockUseCases{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(tt.mock)

			req := httptest.NewRequest("GET", "/api/analytics"+tt.query, nil)
			w := httptest.NewRecorder()

			handler.GetAnalytics(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetAnalytics() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
