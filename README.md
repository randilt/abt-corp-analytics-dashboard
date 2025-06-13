# ABT Analytics Dashboard

High-performance analytics dashboard for ABT Corporation's transaction data. Built with Go backend and React frontend.

Author: [Randil Tharusha](https://randiltharusha.me)

![APP SCREENSHOT](https://github.com/user-attachments/assets/e16777f0-fa59-41e4-8fdb-1c4e727ac3a9)

## Features

- Country-level revenue analysis with transaction counts
- Top 20 most purchased products with stock levels
- Monthly sales volume trends
- Top 30 regions by revenue and items sold
- In-memory and file-based caching for sub-10s load times
- Concurrent data processing for optimal performance

## Tech Stack

- Backend: Go 1.23
- Frontend: React + Vite
- Data Processing: Concurrent Go routines
- Caching: In-memory + File-based
- Containerization: Docker + Docker Compose

## Prerequisites

- Go 1.23 or later
- Node.js 18 or later
- npm or yarn
- Make (optional, for using make commands)

## Dataset Setup

The repository includes dummy data files for reference:

- `data/raw/transactions.csv` (100 rows - 99 records + 1 header)
- `data/processed/cache.json` (cached analytics for the dummy data)

These files are for reference only. Before running the application:

1. Remove the dummy files:

```bash
rm data/raw/transactions.csv
rm data/processed/cache.json
```

2. Place your full transaction dataset in `data/raw/transactions.csv`

Note: The full dataset is large (500+ MB) and contains millions of transaction records. Make sure you have enough disk space and memory to process it.

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

### CSV Processing Configuration

```bash
CSV_FILE_PATH=./data/raw/transactions.csv  # Path to CSV file
CSV_BATCH_SIZE=10000          # Number of records to process in each batch
CSV_WORKER_POOL=8             # Number of concurrent workers (reduce if high resource usage)
CSV_BUFFER_SIZE=65536         # Buffer size for CSV reading
```

### Cache Configuration

```bash
CACHE_FILE_PATH=./data/processed/analytics_cache.json  # Path to cache file
CACHE_TTL=24h                 # Cache time-to-live
```

### Logging Configuration

```bash
LOG_LEVEL=info               # Log level (debug, info, warn, error)
```

Example usage:

```bash
# Run with custom configuration
CSV_BATCH_SIZE=5000 CSV_WORKER_POOL=4 ./bin/server
```

## Accessing the Dashboard

- Frontend: http://localhost:4173
- Backend API: http://localhost:8080/api/v1

## Data Processing

- CSV data is processed concurrently using Go routines
- Results are cached in memory and persisted to disk in a json file (I'm converting the CSV data to JSON as it is flexible and easy to work with)
- Initial data is loaded from `data/raw/transactions.csv`
- Cache TTL: 24 hours
- Data can be refreshed from `/api/v1/analytics/refresh` endpoint

## API Endpoints

- `GET /api/v1/analytics/summary` - Get all analytics data
- `GET /api/v1/analytics/countries` - Country revenue data
- `GET /api/v1/analytics/products` - Top products
- `GET /api/v1/analytics/regions` - Top regions
- `GET /api/v1/analytics/sales` - Monthly sales
- `POST /api/v1/analytics/refresh` - Force data refresh

## Performance

- Fresh data load (No cache): 8-10s
- Subsequent loads (Cache hit): avg 5ms when loading from memory
- Concurrent processing of 5 million records
- Memory usage: ~500mb for in memory cache
- Concurrent processing of:
  - Country revenue
  - Top products
  - Monthly sales
  - Region analysis

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
   - Try refreshing the cache using the refresh endpoint

2. If the frontend build fails:
   - Clear node_modules and reinstall: `rm -rf node_modules && npm install`
   - Ensure you're using Node.js 18 or later
