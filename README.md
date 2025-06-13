# ABT Analytics Dashboard

High-performance analytics dashboard for ABT Corporation's transaction data. Built with Go backend and React frontend.

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

## Setup

1. Clone and build:

```bash
git clone https://github.com/randilt/abt-corp-analytics-dashboard.git
cd abt-corp-analytics-dashboard
make build
```

2. Run with Docker:

Pull from dockerhub:

```bash
docker pull randilt/abt-corp-analytics-dashboard:latest
```

Or start with docker-compose:

```bash
docker-compose up -d
```

3. Access dashboard at `http://localhost:8080`

## Data Processing

- CSV data is processed concurrently using Go routines
- Results are cached in memory and persisted to disk in a json file (I'm converting the CSV data to JSON as it is flexible and easy to work with)
- Initial data is loaded from `data/transactions.csv`
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
