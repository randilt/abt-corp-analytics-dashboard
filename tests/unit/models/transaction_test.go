package models_test

import (
	"testing"
	"time"

	"analytics-dashboard-api/internal/models"
)

func TestTransaction_ParseCSVRow(t *testing.T) {
	tests := []struct {
		name    string
		row     []string
		wantErr bool
		want    models.Transaction
	}{
		{
			name: "valid complete row",
			row: []string{
				"T123", "2023-01-15", "U456", "USA", "California",
				"P789", "Test Product", "Electronics", "299.99", "2",
				"599.98", "100", "2022-12-01",
			},
			wantErr: false,
			want: models.Transaction{
				TransactionID:   "T123",
				TransactionDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
				UserID:          "U456",
				Country:         "USA",
				Region:          "California",
				ProductID:       "P789",
				ProductName:     "Test Product",
				Category:        "Electronics",
				Price:           299.99,
				Quantity:        2,
				TotalPrice:      599.98,
				StockQuantity:   100,
				AddedDate:       time.Date(2022, 12, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "valid row without added_date",
			row: []string{
				"T124", "2023-01-16", "U457", "Canada", "Ontario",
				"P790", "Test Product 2", "Books", "29.99", "1",
				"29.99", "50",
			},
			wantErr: false,
			want: models.Transaction{
				TransactionID:   "T124",
				TransactionDate: time.Date(2023, 1, 16, 0, 0, 0, 0, time.UTC),
				UserID:          "U457",
				Country:         "Canada",
				Region:          "Ontario",
				ProductID:       "P790",
				ProductName:     "Test Product 2",
				Category:        "Books",
				Price:           29.99,
				Quantity:        1,
				TotalPrice:      29.99,
				StockQuantity:   50,
			},
		},
		{
			name:    "insufficient columns",
			row:     []string{"T123", "2023-01-15"},
			wantErr: true,
		},
		{
			name:    "empty transaction_id",
			row:     []string{"", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "2", "599.98", "100"},
			wantErr: true,
		},
		{
			name:    "invalid date format",
			row:     []string{"T123", "invalid-date", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "2", "599.98", "100"},
			wantErr: true,
		},
		{
			name:    "invalid price",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "invalid", "2", "599.98", "100"},
			wantErr: true,
		},
		{
			name:    "negative price",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "-10.00", "2", "599.98", "100"},
			wantErr: true,
		},
		{
			name:    "invalid quantity",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "invalid", "599.98", "100"},
			wantErr: true,
		},
		{
			name:    "zero quantity",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "0", "599.98", "100"},
			wantErr: true,
		},
		{
			name:    "invalid total_price",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "2", "invalid", "100"},
			wantErr: true,
		},
		{
			name:    "invalid stock_quantity",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "2", "599.98", "invalid"},
			wantErr: true,
		},
		{
			name:    "negative stock_quantity",
			row:     []string{"T123", "2023-01-15", "U456", "USA", "California", "P789", "Test Product", "Electronics", "299.99", "2", "599.98", "-10"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var transaction models.Transaction
			err := transaction.ParseCSVRow(tt.row)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCSVRow() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseCSVRow() unexpected error: %v", err)
				return
			}

			// Compare individual fields
			if transaction.TransactionID != tt.want.TransactionID {
				t.Errorf("TransactionID = %v, want %v", transaction.TransactionID, tt.want.TransactionID)
			}
			if !transaction.TransactionDate.Equal(tt.want.TransactionDate) {
				t.Errorf("TransactionDate = %v, want %v", transaction.TransactionDate, tt.want.TransactionDate)
			}
			if transaction.UserID != tt.want.UserID {
				t.Errorf("UserID = %v, want %v", transaction.UserID, tt.want.UserID)
			}
			if transaction.Country != tt.want.Country {
				t.Errorf("Country = %v, want %v", transaction.Country, tt.want.Country)
			}
			if transaction.Region != tt.want.Region {
				t.Errorf("Region = %v, want %v", transaction.Region, tt.want.Region)
			}
			if transaction.Price != tt.want.Price {
				t.Errorf("Price = %v, want %v", transaction.Price, tt.want.Price)
			}
			if transaction.Quantity != tt.want.Quantity {
				t.Errorf("Quantity = %v, want %v", transaction.Quantity, tt.want.Quantity)
			}
			if transaction.TotalPrice != tt.want.TotalPrice {
				t.Errorf("TotalPrice = %v, want %v", transaction.TotalPrice, tt.want.TotalPrice)
			}
			if transaction.StockQuantity != tt.want.StockQuantity {
				t.Errorf("StockQuantity = %v, want %v", transaction.StockQuantity, tt.want.StockQuantity)
			}
		})
	}
}

func TestTransaction_GetMonth(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "january 2023",
			date: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			want: "2023-01",
		},
		{
			name: "december 2022",
			date: time.Date(2022, 12, 31, 23, 59, 59, 0, time.UTC),
			want: "2022-12",
		},
		{
			name: "february 2024",
			date: time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			want: "2024-02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transaction := models.Transaction{
				TransactionDate: tt.date,
			}

			got := transaction.GetMonth()
			if got != tt.want {
				t.Errorf("GetMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_ParseCSVRow_AlternativeDateFormats(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "format 2006-01-02",
			dateStr:  "2023-01-15",
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "format 01/02/2006",
			dateStr:  "01/15/2023",
			expected: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "format 2006-01-02 15:04:05",
			dateStr:  "2023-01-15 14:30:45",
			expected: time.Date(2023, 1, 15, 14, 30, 45, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:    "invalid format",
			dateStr: "15-01-2023",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row := []string{
				"T123", tt.dateStr, "U456", "USA", "California",
				"P789", "Test Product", "Electronics", "299.99", "2",
				"599.98", "100", "2022-12-01",
			}

			var transaction models.Transaction
			err := transaction.ParseCSVRow(row)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCSVRow() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseCSVRow() unexpected error: %v", err)
				return
			}

			if !transaction.TransactionDate.Equal(tt.expected) {
				t.Errorf("TransactionDate = %v, want %v", transaction.TransactionDate, tt.expected)
			}
		})
	}
}
