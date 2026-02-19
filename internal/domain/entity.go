package domain

import (
	"errors"
	"time"
)

// Item представляет финансовую транзакцию или запись
type Item struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`      // "income" или "expense"
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Analytics представляет агрегированную аналитику
type Analytics struct {
	Sum        float64 `json:"sum"`
	Avg        float64 `json:"avg"`
	Count      int64   `json:"count"`
	Median     float64 `json:"median"`
	Percentile float64 `json:"percentile_90"`
}

// Validate проверяет корректность данных
func (i *Item) Validate() error {
	if i.Amount < 0 {
		return errors.New("amount cannot be negative")
	}
	if i.Type != "income" && i.Type != "expense" {
		return errors.New("type must be 'income' or 'expense'")
	}
	if i.Category == "" {
		return errors.New("category is required")
	}
	if i.Date.IsZero() {
		return errors.New("date is required")
	}
	return nil
}
