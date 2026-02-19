package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/yourusername/analytics-service/internal/domain"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN not set, skipping integration tests")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id BIGSERIAL PRIMARY KEY,
			type VARCHAR(20) NOT NULL CHECK (type IN ('income', 'expense')),
			amount DECIMAL(15, 2) NOT NULL CHECK (amount >= 0),
			category VARCHAR(100) NOT NULL,
			date TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	cleanup := func() {
		db.Exec("DROP TABLE IF EXISTS items")
		db.Close()
	}

	return db, cleanup
}

func TestRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	item := &domain.Item{
		Type:      "income",
		Amount:    1000.50,
		Category:  "Salary",
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := repo.Create(ctx, item)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if item.ID == 0 {
		t.Error("Create() did not set ID")
	}
}

func TestRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	// Create test item
	item := &domain.Item{
		Type:      "expense",
		Amount:    500.00,
		Category:  "Food",
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(ctx, item)

	// Get item
	got, err := repo.GetByID(ctx, item.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != item.ID {
		t.Errorf("GetByID() ID = %v, want %v", got.ID, item.ID)
	}
	if got.Type != item.Type {
		t.Errorf("GetByID() Type = %v, want %v", got.Type, item.Type)
	}
	if got.Amount != item.Amount {
		t.Errorf("GetByID() Amount = %v, want %v", got.Amount, item.Amount)
	}
}

func TestRepository_GetAll(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	// Create test items
	now := time.Now()
	items := []*domain.Item{
		{Type: "income", Amount: 1000.00, Category: "Salary", Date: now, CreatedAt: now, UpdatedAt: now},
		{Type: "expense", Amount: 500.00, Category: "Food", Date: now.AddDate(0, 0, -1), CreatedAt: now, UpdatedAt: now},
		{Type: "income", Amount: 2000.00, Category: "Bonus", Date: now.AddDate(0, 0, -2), CreatedAt: now, UpdatedAt: now},
	}

	for _, item := range items {
		repo.Create(ctx, item)
	}

	// Get all items
	got, err := repo.GetAll(ctx, nil, nil)
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(got) != len(items) {
		t.Errorf("GetAll() returned %d items, want %d", len(got), len(items))
	}
}

func TestRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	// Create test item
	item := &domain.Item{
		Type:      "income",
		Amount:    1000.00,
		Category:  "Salary",
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(ctx, item)

	// Update item
	item.Amount = 1500.00
	item.Category = "Salary + Bonus"
	item.UpdatedAt = time.Now()

	err := repo.Update(ctx, item)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	got, _ := repo.GetByID(ctx, item.ID)
	if got.Amount != 1500.00 {
		t.Errorf("Update() Amount = %v, want %v", got.Amount, 1500.00)
	}
	if got.Category != "Salary + Bonus" {
		t.Errorf("Update() Category = %v, want %v", got.Category, "Salary + Bonus")
	}
}

func TestRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	// Create test item
	item := &domain.Item{
		Type:      "expense",
		Amount:    300.00,
		Category:  "Transport",
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Create(ctx, item)

	// Delete item
	err := repo.Delete(ctx, item.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, item.ID)
	if err == nil {
		t.Error("Delete() item still exists")
	}
}

func TestRepository_GetAnalytics(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	// Create test items
	now := time.Now()
	amounts := []float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}
	
	for _, amount := range amounts {
		item := &domain.Item{
			Type:      "income",
			Amount:    amount,
			Category:  "Test",
			Date:      now,
			CreatedAt: now,
			UpdatedAt: now,
		}
		repo.Create(ctx, item)
	}

	// Get analytics
	from := now.AddDate(0, 0, -1)
	to := now.AddDate(0, 0, 1)
	
	analytics, err := repo.GetAnalytics(ctx, from, to)
	if err != nil {
		t.Fatalf("GetAnalytics() error = %v", err)
	}

	// Verify results
	expectedSum := 5500.0
	if analytics.Sum != expectedSum {
		t.Errorf("GetAnalytics() Sum = %v, want %v", analytics.Sum, expectedSum)
	}

	expectedAvg := 550.0
	if analytics.Avg != expectedAvg {
		t.Errorf("GetAnalytics() Avg = %v, want %v", analytics.Avg, expectedAvg)
	}

	if analytics.Count != 10 {
		t.Errorf("GetAnalytics() Count = %v, want %v", analytics.Count, 10)
	}

	expectedMedian := 550.0
	if analytics.Median != expectedMedian {
		t.Errorf("GetAnalytics() Median = %v, want %v", analytics.Median, expectedMedian)
	}

	// 90th percentile should be 900
	expectedPercentile := 900.0
	if analytics.Percentile != expectedPercentile {
		t.Errorf("GetAnalytics() Percentile = %v, want %v", analytics.Percentile, expectedPercentile)
	}
}

func TestRepository_GetAnalytics_EmptyData(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := &repository{db: db}
	ctx := context.Background()

	// Get analytics with no data
	from := time.Now().AddDate(0, 0, -1)
	to := time.Now()
	
	analytics, err := repo.GetAnalytics(ctx, from, to)
	if err != nil {
		t.Fatalf("GetAnalytics() error = %v", err)
	}

	// All values should be 0
	if analytics.Sum != 0 || analytics.Avg != 0 || analytics.Count != 0 {
		t.Errorf("GetAnalytics() with empty data should return zeros, got %+v", analytics)
	}
}


