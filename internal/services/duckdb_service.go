package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"analytics-dashboard-api/internal/models"
	"analytics-dashboard-api/pkg/logger"

	_ "github.com/marcboeker/go-duckdb"
)

type DuckDBService struct {
	db     *sql.DB
	logger logger.Logger
}

func NewDuckDBService(logger logger.Logger) (*DuckDBService, error) {
	// Create in-memory DuckDB database
	db, err := sql.Open("duckdb", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open DuckDB: %w", err)
	}

	service := &DuckDBService{
		db:     db,
		logger: logger,
	}

	// Create transactions table
	if err := service.createTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return service, nil
}

func (s *DuckDBService) Close() error {
	return s.db.Close()
}

func (s *DuckDBService) createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		transaction_id VARCHAR,
		transaction_date DATE,
		user_id VARCHAR,
		country VARCHAR,
		region VARCHAR,
		product_id VARCHAR,
		product_name VARCHAR,
		category VARCHAR,
		price DECIMAL(10,2),
		quantity INTEGER,
		total_price DECIMAL(10,2),
		stock_quantity INTEGER,
		added_date DATE
	)`
	
	_, err := s.db.Exec(createTableSQL)
	return err
}

func (s *DuckDBService) LoadFromCSV(csvPath string) error {
	startTime := time.Now()
	s.logger.Info("Loading CSV data into DuckDB", "file", csvPath)

	// Use DuckDB's CSV reader to load data directly
	loadSQL := fmt.Sprintf(`
		INSERT INTO transactions 
		SELECT 
			transaction_id,
			CAST(transaction_date AS DATE) as transaction_date,
			user_id,
			country,
			region,
			product_id,
			product_name,
			category,
			CAST(price AS DECIMAL(10,2)) as price,
			CAST(quantity AS INTEGER) as quantity,
			CAST(total_price AS DECIMAL(10,2)) as total_price,
			CAST(stock_quantity AS INTEGER) as stock_quantity,
			CAST(added_date AS DATE) as added_date
		FROM read_csv_auto('%s', header=true)
	`, csvPath)

	_, err := s.db.Exec(loadSQL)
	if err != nil {
		return fmt.Errorf("failed to load CSV: %w", err)
	}

	// Get row count
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to get row count: %w", err)
	}

	s.logger.Info("CSV data loaded successfully", 
		"records", count, 
		"duration", time.Since(startTime))

	return nil
}

func (s *DuckDBService) GetCountryRevenue(ctx context.Context, limit, offset int) ([]models.CountryRevenue, error) {
	query := `
		SELECT 
			country,
			product_name,
			CAST(SUM(total_price) AS DOUBLE) as total_revenue,
			COUNT(*) as transaction_count
		FROM transactions 
		GROUP BY country, product_name
		ORDER BY total_revenue DESC
		LIMIT ? OFFSET ?
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query country revenue: %w", err)
	}
	defer rows.Close()

	var results []models.CountryRevenue
	for rows.Next() {
		var cr models.CountryRevenue
		err := rows.Scan(
			&cr.Country,
			&cr.ProductName,
			&cr.TotalRevenue,
			&cr.TransactionCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan country revenue: %w", err)
		}
		results = append(results, cr)
	}

	return results, nil
}

func (s *DuckDBService) GetTopProducts(ctx context.Context) ([]models.ProductFrequency, error) {
	query := `
		SELECT 
			product_id,
			product_name,
			SUM(quantity) as purchase_count,
			MAX(stock_quantity) as stock_quantity
		FROM transactions 
		GROUP BY product_id, product_name
		ORDER BY purchase_count DESC
		LIMIT 20
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query top products: %w", err)
	}
	defer rows.Close()

	var results []models.ProductFrequency
	for rows.Next() {
		var pf models.ProductFrequency
		err := rows.Scan(
			&pf.ProductID,
			&pf.ProductName,
			&pf.PurchaseCount,
			&pf.StockQuantity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top products: %w", err)
		}
		results = append(results, pf)
	}

	return results, nil
}

func (s *DuckDBService) GetMonthlySales(ctx context.Context) ([]models.MonthlySales, error) {
	query := `
		SELECT 
			STRFTIME('%Y-%m', transaction_date) as month,
			CAST(SUM(total_price) AS DOUBLE) as sales_volume,
			SUM(quantity) as item_count
		FROM transactions 
		GROUP BY STRFTIME('%Y-%m', transaction_date)
		ORDER BY month
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly sales: %w", err)
	}
	defer rows.Close()

	var results []models.MonthlySales
	for rows.Next() {
		var ms models.MonthlySales
		err := rows.Scan(
			&ms.Month,
			&ms.SalesVolume,
			&ms.ItemCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan monthly sales: %w", err)
		}
		results = append(results, ms)
	}

	return results, nil
}

func (s *DuckDBService) GetTopRegions(ctx context.Context) ([]models.RegionRevenue, error) {
	query := `
		SELECT 
			region,
			CAST(SUM(total_price) AS DOUBLE) as total_revenue,
			SUM(quantity) as items_sold
		FROM transactions 
		GROUP BY region
		ORDER BY total_revenue DESC
		LIMIT 30
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query top regions: %w", err)
	}
	defer rows.Close()

	var results []models.RegionRevenue
	for rows.Next() {
		var rr models.RegionRevenue
		err := rows.Scan(
			&rr.Region,
			&rr.TotalRevenue,
			&rr.ItemsSold,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top regions: %w", err)
		}
		results = append(results, rr)
	}

	return results, nil
}

func (s *DuckDBService) GetTotalRecords(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	return count, err
}

func (s *DuckDBService) GetCountryRevenueCount(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM (
			SELECT DISTINCT country, product_name 
			FROM transactions
		)
	`).Scan(&count)
	return count, err
}
