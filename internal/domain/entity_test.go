package domain

import (
	"testing"
	"time"
)

func TestItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid income item",
			item: Item{
				Type:     "income",
				Amount:   1000.50,
				Category: "Salary",
				Date:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid expense item",
			item: Item{
				Type:     "expense",
				Amount:   500.00,
				Category: "Food",
				Date:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "negative amount",
			item: Item{
				Type:     "income",
				Amount:   -100.00,
				Category: "Salary",
				Date:     time.Now(),
			},
			wantErr: true,
			errMsg:  "amount cannot be negative",
		},
		{
			name: "invalid type",
			item: Item{
				Type:     "invalid",
				Amount:   100.00,
				Category: "Test",
				Date:     time.Now(),
			},
			wantErr: true,
			errMsg:  "type must be 'income' or 'expense'",
		},
		{
			name: "empty category",
			item: Item{
				Type:     "income",
				Amount:   100.00,
				Category: "",
				Date:     time.Now(),
			},
			wantErr: true,
			errMsg:  "category is required",
		},
		{
			name: "zero date",
			item: Item{
				Type:     "income",
				Amount:   100.00,
				Category: "Test",
				Date:     time.Time{},
			},
			wantErr: true,
			errMsg:  "date is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Item.Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
