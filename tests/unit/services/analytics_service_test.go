package services_test

import (
	"fmt"
	"testing"
	"time"

	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/internal/services"
)

// Mock logger for testing
type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...interface{})  {}
func (m *mockLogger) Error(msg string, args ...interface{}) {}
func (m *mockLogger) Debug(msg string, args ...interface{}) {}
func (m *mockLogger) Warn(msg string, args ...interface{})  {}

func createTestTransactions() []models.Transaction {
	return []models.Transaction{
		{
			TransactionID:   "T1",
			TransactionDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			UserID:          "U1",
			Country:         "USA",
			Region:          "California",
			ProductID:       "P1",
			ProductName:     "Product A",
			Category:        "Electronics",
			Price:           100.0,
			Quantity:        2,
			TotalPrice:      200.0,
			StockQuantity:   50,
		},
		{
			TransactionID:   "T2",
			TransactionDate: time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC),
			UserID:          "U2",
			Country:         "USA",
			Region:          "Texas",
			ProductID:       "P2",
			ProductName:     "Product B",
			Category:        "Books",
			Price:           25.0,
			Quantity:        3,
			TotalPrice:      75.0,
			StockQuantity:   30,
		},
		{
			TransactionID:   "T3",
			TransactionDate: time.Date(2023, 2, 10, 0, 0, 0, 0, time.UTC),
			UserID:          "U3",
			Country:         "Canada",
			Region:          "Ontario",
			ProductID:       "P1",
			ProductName:     "Product A",
			Category:        "Electronics",
			Price:           100.0,
			Quantity:        1,
			TotalPrice:      100.0,
			StockQuantity:   50,
		},
		{
			TransactionID:   "T4",
			TransactionDate: time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			UserID:          "U4",
			Country:         "USA",
			Region:          "California",
			ProductID:       "P3",
			ProductName:     "Product C",
			Category:        "Clothing",
			Price:           50.0,
			Quantity:        4,
			TotalPrice:      200.0,
			StockQuantity:   20,
		},
		{
			TransactionID:   "T5",
			TransactionDate: time.Date(2023, 3, 5, 0, 0, 0, 0, time.UTC),
			UserID:          "U5",
			Country:         "Germany",
			Region:          "Bavaria",
			ProductID:       "P2",
			ProductName:     "Product B",
			Category:        "Books",
			Price:           25.0,
			Quantity:        5,
			TotalPrice:      125.0,
			StockQuantity:   30,
		},
	}
}

func TestAnalyticsService_GenerateAnalytics(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := []models.Transaction{
		{
			TransactionID:   "T1",
			TransactionDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UserID:          "U1",
			ProductID:       "P1",
			ProductName:     "Product A",
			Category:        "Electronics",
			Price:           100.0,
			Quantity:        2,
			TotalPrice:      200.0,
			Country:         "USA",
			Region:          "California",
			StockQuantity:   50,
			AddedDate:       time.Now(),
		},
	}

	// Add a small delay to ensure processing time is measurable
	time.Sleep(10 * time.Millisecond)
	result := service.GenerateAnalytics(transactions)

	if result == nil {
		t.Fatal("GenerateAnalytics() returned nil")
	}

	if result.ProcessingTimeMs <= 0 {
		t.Errorf("ProcessingTimeMs should be positive, got %d", result.ProcessingTimeMs)
	}

	if result.TotalRecords != 1 {
		t.Errorf("TotalRecords = %d, want 1", result.TotalRecords)
	}

	if len(result.CountryRevenue) != 1 {
		t.Errorf("CountryRevenue length = %d, want 1", len(result.CountryRevenue))
	}

	if len(result.TopProducts) != 1 {
		t.Errorf("TopProducts length = %d, want 1", len(result.TopProducts))
	}

	if len(result.MonthlySales) != 1 {
		t.Errorf("MonthlySales length = %d, want 1", len(result.MonthlySales))
	}

	if len(result.TopRegions) != 1 {
		t.Errorf("TopRegions length = %d, want 1", len(result.TopRegions))
	}

	// Verify processing time is reasonable (should be less than 100ms for a single record)
	if result.ProcessingTimeMs > 100 {
		t.Errorf("Processing time too high: %d ms", result.ProcessingTimeMs)
	}
}

func TestAnalyticsService_GenerateCountryRevenue(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := createTestTransactions()

	result := service.GenerateAnalytics(transactions)
	countryRevenue := result.CountryRevenue

	// Should have multiple country-product combinations
	if len(countryRevenue) < 3 {
		t.Errorf("Expected at least 3 country-product combinations, got %d", len(countryRevenue))
	}

	// Check sorting (should be by total revenue descending)
	for i := 1; i < len(countryRevenue); i++ {
		if countryRevenue[i].TotalRevenue > countryRevenue[i-1].TotalRevenue {
			t.Errorf("CountryRevenue not sorted correctly: %f > %f at positions %d, %d",
				countryRevenue[i].TotalRevenue, countryRevenue[i-1].TotalRevenue, i, i-1)
		}
	}

	// Verify specific data
	found := false
	for _, cr := range countryRevenue {
		if cr.Country == "USA" && cr.ProductName == "Product A" {
			if cr.TotalRevenue != 200.0 {
				t.Errorf("USA Product A revenue = %f, want 200.0", cr.TotalRevenue)
			}
			if cr.TransactionCount != 1 {
				t.Errorf("USA Product A transaction count = %d, want 1", cr.TransactionCount)
			}
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find USA Product A in country revenue")
	}
}

func TestAnalyticsService_GenerateTopProducts(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := createTestTransactions()

	result := service.GenerateAnalytics(transactions)
	topProducts := result.TopProducts

	// Should have 3 unique products
	if len(topProducts) != 3 {
		t.Errorf("Expected 3 unique products, got %d", len(topProducts))
	}

	// Check sorting (should be by purchase count descending)
	for i := 1; i < len(topProducts); i++ {
		if topProducts[i].PurchaseCount > topProducts[i-1].PurchaseCount {
			t.Errorf("TopProducts not sorted correctly: %d > %d at positions %d, %d",
				topProducts[i].PurchaseCount, topProducts[i-1].PurchaseCount, i, i-1)
		}
	}

	// Verify Product B has highest purchase count (3+5=8)
	if topProducts[0].ProductID != "P2" {
		t.Errorf("Expected P2 to be top product, got %s", topProducts[0].ProductID)
	}
	if topProducts[0].PurchaseCount != 8 {
		t.Errorf("Expected P2 purchase count to be 8, got %d", topProducts[0].PurchaseCount)
	}
}

func TestAnalyticsService_GenerateMonthlySales(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := createTestTransactions()

	result := service.GenerateAnalytics(transactions)
	monthlySales := result.MonthlySales

	// Should have 3 months (2023-01, 2023-02, 2023-03)
	if len(monthlySales) != 3 {
		t.Errorf("Expected 3 months, got %d", len(monthlySales))
	}

	// Check sorting (should be by month ascending)
	expectedMonths := []string{"2023-01", "2023-02", "2023-03"}
	for i, expected := range expectedMonths {
		if monthlySales[i].Month != expected {
			t.Errorf("Month at position %d = %s, want %s", i, monthlySales[i].Month, expected)
		}
	}

	// Verify January sales (T1: 200 + T2: 75 = 275)
	january := monthlySales[0]
	if january.SalesVolume != 275.0 {
		t.Errorf("January sales volume = %f, want 275.0", january.SalesVolume)
	}
	if january.ItemCount != 5 { // 2+3
		t.Errorf("January item count = %d, want 5", january.ItemCount)
	}
}

func TestAnalyticsService_GenerateTopRegions(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := createTestTransactions()

	result := service.GenerateAnalytics(transactions)
	topRegions := result.TopRegions

	// Should have 4 unique regions
	if len(topRegions) != 4 {
		t.Errorf("Expected 4 unique regions, got %d", len(topRegions))
	}

	// Check sorting (should be by total revenue descending)
	for i := 1; i < len(topRegions); i++ {
		if topRegions[i].TotalRevenue > topRegions[i-1].TotalRevenue {
			t.Errorf("TopRegions not sorted correctly: %f > %f at positions %d, %d",
				topRegions[i].TotalRevenue, topRegions[i-1].TotalRevenue, i, i-1)
		}
	}

	// Verify California has highest revenue (T1: 200 + T4: 200 = 400)
	california := topRegions[0]
	if california.Region != "California" {
		t.Errorf("Expected California to be top region, got %s", california.Region)
	}
	if california.TotalRevenue != 400.0 {
		t.Errorf("California revenue = %f, want 400.0", california.TotalRevenue)
	}
	if california.ItemsSold != 6 { // 2+4
		t.Errorf("California items sold = %d, want 6", california.ItemsSold)
	}
}

func TestAnalyticsService_EmptyTransactions(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := []models.Transaction{}

	result := service.GenerateAnalytics(transactions)

	if result.TotalRecords != 0 {
		t.Errorf("TotalRecords = %d, want 0", result.TotalRecords)
	}

	if len(result.CountryRevenue) != 0 {
		t.Errorf("CountryRevenue should be empty, got %d items", len(result.CountryRevenue))
	}

	if len(result.TopProducts) != 0 {
		t.Errorf("TopProducts should be empty, got %d items", len(result.TopProducts))
	}

	if len(result.MonthlySales) != 0 {
		t.Errorf("MonthlySales should be empty, got %d items", len(result.MonthlySales))
	}

	if len(result.TopRegions) != 0 {
		t.Errorf("TopRegions should be empty, got %d items", len(result.TopRegions))
	}
}

func TestAnalyticsService_TopProductsLimit(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)

	// Create 25 different products to test the limit of 20
	transactions := make([]models.Transaction, 25)
	for i := 0; i < 25; i++ {
		transactions[i] = models.Transaction{
			TransactionID:   fmt.Sprintf("T%d", i+1),
			TransactionDate: time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			UserID:          fmt.Sprintf("U%d", i+1),
			Country:         "USA",
			Region:          "California",
			ProductID:       fmt.Sprintf("P%d", i+1),
			ProductName:     fmt.Sprintf("Product %d", i+1),
			Category:        "Electronics",
			Price:           100.0,
			Quantity:        i + 1, // Different quantities to ensure different ranking
			TotalPrice:      100.0 * float64(i+1),
			StockQuantity:   50,
		}
	}

	result := service.GenerateAnalytics(transactions)

	// Should limit to top 20 products
	if len(result.TopProducts) != 20 {
		t.Errorf("Expected top 20 products, got %d", len(result.TopProducts))
	}

	// Verify the top product has the highest quantity (25)
	if result.TopProducts[0].PurchaseCount != 25 {
		t.Errorf("Top product purchase count = %d, want 25", result.TopProducts[0].PurchaseCount)
	}
}

func TestAnalyticsService_TopRegionsLimit(t *testing.T) {
	logger := &mockLogger{}
	service := services.NewAnalyticsService(logger)
	transactions := make([]models.Transaction, 40)

	// Create 40 transactions with different regions
	for i := 0; i < 40; i++ {
		transactions[i] = models.Transaction{
			TransactionID:   fmt.Sprintf("T%d", i+1),
			TransactionDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UserID:          fmt.Sprintf("U%d", i+1),
			ProductID:       fmt.Sprintf("P%d", i+1),
			ProductName:     fmt.Sprintf("Product %d", i+1),
			Category:        "Electronics",
			Price:           100.0,
			Quantity:        1,
			TotalPrice:      100.0,
			Country:         "USA",
			Region:          fmt.Sprintf("Region %d", i+1),
			StockQuantity:   50,
			AddedDate:       time.Now(),
		}
	}

	result := service.GenerateAnalytics(transactions)

	if len(result.TopRegions) != 30 {
		t.Errorf("Expected top 30 regions, got %d", len(result.TopRegions))
	}

	// Verify regions are sorted by revenue
	for i := 1; i < len(result.TopRegions); i++ {
		if result.TopRegions[i-1].TotalRevenue < result.TopRegions[i].TotalRevenue {
			t.Errorf("Regions not sorted by revenue: %v < %v",
				result.TopRegions[i-1].TotalRevenue,
				result.TopRegions[i].TotalRevenue)
		}
	}
}
