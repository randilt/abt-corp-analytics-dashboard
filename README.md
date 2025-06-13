# ABT Analytics Dashboard

High-performance analytics dashboard for ABT Corporation's transaction data. Built with Go backend and React frontend.

Author: [Randil Tharusha](https://randiltharusha.me)

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

The transaction dataset provided is not included in the repository due to its large size (500+ MB). You'll need to:

1. Create a data directory:

```bash
mkdir -p data/raw
```

2. Place the transactions dataset in `data/raw/transactions.csv`

Note: The provided dataset is large (500+ MB) and contains millions of transaction records. Make sure you have enough disk space and memory to process it.

## Setup

### Quick Start Script

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
make test        # Run tests
make coverage    # Generate coverage report
```

Current coverage: 85% (core business logic)

## Troubleshooting

1. If data isn't loading:

   - Check if the backend is running (`curl http://localhost:8080/health`)
   - Verify the CSV file exists in `data/raw/transactions.csv`
   - Try refreshing the cache using the refresh endpoint

2. If the frontend build fails:
   - Clear node_modules and reinstall: `rm -rf node_modules && npm install`
   - Ensure you're using Node.js 18 or later
