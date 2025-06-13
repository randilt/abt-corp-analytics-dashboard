// internal/models/transaction.go - CORRECTED VERSION

package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a single transaction record
type Transaction struct {
	TransactionID   string    `json:"transaction_id" csv:"transaction_id"`
	TransactionDate time.Time `json:"transaction_date" csv:"transaction_date"`
	UserID          string    `json:"user_id" csv:"user_id"`
	Country         string    `json:"country" csv:"country"`
	Region          string    `json:"region" csv:"region"`
	ProductID       string    `json:"product_id" csv:"product_id"`
	ProductName     string    `json:"product_name" csv:"product_name"`
	Category        string    `json:"category" csv:"category"`
	Price           float64   `json:"price" csv:"price"`
	Quantity        int       `json:"quantity" csv:"quantity"`
	TotalPrice      float64   `json:"total_price" csv:"total_price"`
	StockQuantity   int       `json:"stock_quantity" csv:"stock_quantity"`
	AddedDate       time.Time `json:"added_date" csv:"added_date"`
}

// ParseCSVRow converts a CSV row to Transaction
func (t *Transaction) ParseCSVRow(row []string) error {
	if len(row) < 12 {
		return fmt.Errorf("insufficient columns: got %d, need at least 12", len(row))
	}

	// Basic field assignment with validation
	t.TransactionID = strings.TrimSpace(row[0])
	if t.TransactionID == "" {
		return fmt.Errorf("empty transaction_id")
	}
	
	// Parse transaction date
	if dateStr := strings.TrimSpace(row[1]); dateStr != "" {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			t.TransactionDate = date
		} else {
			// Try alternative date formats
			if date, err := time.Parse("01/02/2006", dateStr); err == nil {
				t.TransactionDate = date
			} else if date, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
				t.TransactionDate = date
			} else {
				return fmt.Errorf("invalid transaction_date: %s", dateStr)
			}
		}
	}
	
	t.UserID = strings.TrimSpace(row[2])
	t.Country = strings.TrimSpace(row[3])
	t.Region = strings.TrimSpace(row[4])
	t.ProductID = strings.TrimSpace(row[5])
	t.ProductName = strings.TrimSpace(row[6])
	t.Category = strings.TrimSpace(row[7])
	
	// Parse numeric fields with validation
	if priceStr := strings.TrimSpace(row[8]); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil && price >= 0 {
			t.Price = price
		} else {
			return fmt.Errorf("invalid price: %s", priceStr)
		}
	}
	
	if qtyStr := strings.TrimSpace(row[9]); qtyStr != "" {
		if qty, err := strconv.Atoi(qtyStr); err == nil && qty > 0 {
			t.Quantity = qty
		} else {
			return fmt.Errorf("invalid quantity: %s", qtyStr)
		}
	}
	
	if totalStr := strings.TrimSpace(row[10]); totalStr != "" {
		if total, err := strconv.ParseFloat(totalStr, 64); err == nil && total >= 0 {
			t.TotalPrice = total
		} else {
			return fmt.Errorf("invalid total_price: %s", totalStr)
		}
	}
	
	if stockStr := strings.TrimSpace(row[11]); stockStr != "" {
		if stock, err := strconv.Atoi(stockStr); err == nil && stock >= 0 {
			t.StockQuantity = stock
		} else {
			return fmt.Errorf("invalid stock_quantity: %s", stockStr)
		}
	}
	
	// Parse added date if exists
	if len(row) > 12 {
		if dateStr := strings.TrimSpace(row[12]); dateStr != "" {
			if date, err := time.Parse("2006-01-02", dateStr); err == nil {
				t.AddedDate = date
			} else if date, err := time.Parse("01/02/2006", dateStr); err == nil {
				t.AddedDate = date
			}
			// If parsing fails, just leave AddedDate as zero value
		}
	}
	
	return nil
}

// GetMonth returns the month in YYYY-MM format for grouping
func (t *Transaction) GetMonth() string {
	return t.TransactionDate.Format("2006-01")
}