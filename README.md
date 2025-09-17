# ABT Analytics Dashboard

High-performance analytics dashboard for ABT Corporation's transaction data. Built with Go backend, React frontend, and powered by DuckDB for blazing-fast SQL analytics.

Author: [Randil Tharusha](https://randiltharusha.me)

![APP SCREENSHOT](https://github.com/user-attachments/assets/e16777f0-fa59-41e4-8fdb-1c4e727ac3a9)

## Features

- **Country-level revenue analysis** with transaction counts and pagination
- **Top 20 most purchased products** with stock levels
- **Monthly sales volume trends** with date-based aggregation
- **Top 30 regions by revenue** and items sold
- **Blazing-fast SQL queries** using DuckDB's columnar storage
- **Real-time analytics** with automatic CSV loading
- **Concurrent data processing** for optimal performance
- **Memory-efficient** in-memory database with instant queries

## Tech Stack

- **Backend**: Go 1.24
- **Frontend**: React + Vite
- **Database**: DuckDB (in-memory columnar analytics)
- **Data Processing**: SQL-powered queries
- **Containerization**: Docker + Docker Compose

## Prerequisites

- Go 1.24 or later
- Node.js 18 or later
- npm or yarn
- Make (optional, for using make commands)

## Dataset Setup

The repository includes dummy data files for reference:

- `data/raw/transactions.csv` (100 rows - 99 records + 1 header)

These files are for reference only. Before running the application:

1. Remove the dummy files:

```bash
rm data/raw/transactions.csv
```

2. Place your full transaction dataset in `data/raw/transactions.csv`

Note: DuckDB can handle datasets of any size efficiently thanks to its columnar storage and SQL optimization.

## Setup

### Quick Start Script ⚡

Use `start.sh` script that automates the setup and startup process:

```bash
# Make the script executable
chmod +x start.sh

# Run the script
./start.sh
```

The script will:

1. Install backend dependencies
2. Build the backend
3. Install frontend dependencies
4. Build the frontend
5. Start both services

### Manual Setup

1. Clone the repository:

```bash
git clone https://github.com/randilt/abt-corp-analytics-dashboard.git
cd abt-corp-analytics-dashboard
```

2. Backend Setup:

```bash
# Install Go dependencies
go mod download

# Build the backend
go build -o bin/server cmd/server/main.go

# Run the backend
./bin/server
```

3. Frontend Setup:

```bash
cd web
npm install
npm run build
npm run preview
```

### Docker Setup

Pull from dockerhub:

```bash
docker pull randilt/abt-corp-analytics-dashboard:latest
```

Or start with docker-compose:

```bash
docker-compose up -d
```

## Configuration

The application can be configured using environment variables:

### Server Configuration

```bash
SERVER_HOST=localhost          # Server hostname
SERVER_PORT=8080              # Server port
SERVER_READ_TIMEOUT=15s       # Read timeout
SERVER_WRITE_TIMEOUT=15s      # Write timeout
SERVER_IDLE_TIMEOUT=60s       # Idle timeout
```

### CSV Configuration

```bash
CSV_FILE_PATH=./data/raw/transactions.csv  # Path to CSV file
```

### Logging Configuration

```bash
LOG_LEVEL=info               # Log level (debug, info, warn, error)
```

Example usage:

```bash
# Run with custom configuration
CSV_FILE_PATH=./my-data.csv ./bin/server
```

## Accessing the Dashboard

- Frontend: http://localhost:4173
- Backend API: http://localhost:8080/api/v1

## Data Processing

- **DuckDB Integration**: CSV data is loaded directly into DuckDB's in-memory columnar database
- **SQL-Powered Analytics**: All analytics are generated using optimized SQL queries
- **Automatic Loading**: Data is loaded fresh on every application startup
- **Memory Efficient**: Only loads data as needed for each query
- **Real-time**: No caching needed - queries are always fresh and fast

## API Endpoints

- `GET /api/v1/analytics` - Get all analytics data summary
- `GET /api/v1/analytics/stats` - Get analytics statistics
- `GET /api/v1/analytics/country-revenue?limit=100&offset=0` - Country revenue data with pagination
- `GET /api/v1/analytics/top-products` - Top 20 products
- `GET /api/v1/analytics/monthly-sales` - Monthly sales trends
- `GET /api/v1/analytics/top-regions` - Top 30 regions
- `POST /api/v1/analytics/refresh` - Force data reload
- `GET /health` - Health check
- `GET /ready` - Readiness check

## Performance

- **CSV Loading**: ~25ms for 99 records
- **Query Performance**: Sub-millisecond response times
- **Memory Usage**: Minimal - only loads what's needed
- **Scalability**: Handles datasets of any size efficiently
- **Concurrent Queries**: All analytics generated in parallel

## Why DuckDB?

- **Columnar Storage**: Optimized for analytical workloads
- **SQL Interface**: Familiar and powerful query language
- **Memory Efficient**: Only loads necessary data
- **Blazing Fast**: Sub-second queries on large datasets
- **No External Dependencies**: Embedded database, no server needed
- **CSV Native**: Direct CSV loading without preprocessing

## Testing

```bash
# Run all tests
make test

# Generate coverage report
make test-coverage

# Generate detailed coverage report
make test-coverage-detailed
```

The coverage report will be generated in `./coverage/coverage.html`. Open this file in your browser to view:

- Line-by-line coverage
- Function coverage
- Overall coverage statistics

Note: Coverage report is generated even if some tests fail.

## Troubleshooting

1. If data isn't loading:

   - Check if the backend is running (`curl http://localhost:8080/health`)
   - Verify the CSV file exists in `data/raw/transactions.csv`
   - Check server logs for DuckDB initialization errors

2. If the frontend build fails:

   - Clear node_modules and reinstall: `rm -rf node_modules && npm install`
   - Ensure you're using Node.js 18 or later

3. If queries are slow:
   - Check CSV file format and column names
   - Verify data types are correct (dates, numbers)
   - Check server logs for SQL errors

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   React Frontend │    │   Go Backend     │    │   DuckDB        │
│                 │    │                  │    │                 │
│ - Dashboard UI  │◄──►│ - REST API       │◄──►│ - In-memory DB  │
│ - Charts        │    │ - Handlers       │    │ - SQL Queries   │
│ - Real-time     │    │ - Middleware     │    │ - Columnar      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌──────────────────┐
                       │   CSV File       │
                       │                  │
                       │ - transactions   │
                       │ - Auto-loaded    │
                       └──────────────────┘
```
