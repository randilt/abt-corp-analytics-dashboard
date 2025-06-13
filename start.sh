#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}[*] $1${NC}"
}

# Function to print error messages
print_error() {
    echo -e "${RED}[!] $1${NC}"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.23 or later."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    print_error "Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    print_error "npm is not installed. Please install npm."
    exit 1
fi

# Create bin directory if it doesn't exist
mkdir -p bin

# Backend setup
print_status "Setting up backend..."
go mod download
if [ $? -ne 0 ]; then
    print_error "Failed to download Go dependencies"
    exit 1
fi

print_status "Building backend..."
go build -o bin/server cmd/server/main.go
if [ $? -ne 0 ]; then
    print_error "Failed to build backend"
    exit 1
fi

# Frontend setup
print_status "Setting up frontend..."
cd web
npm install
if [ $? -ne 0 ]; then
    print_error "Failed to install frontend dependencies"
    exit 1
fi

print_status "Building frontend..."
npm run build
if [ $? -ne 0 ]; then
    print_error "Failed to build frontend"
    exit 1
fi

# Start services
print_status "Starting services..."

# Start backend in background
cd ..
./bin/server &
BACKEND_PID=$!

# Wait for backend to be ready
print_status "Waiting for backend to be ready..."
max_attempts=30
attempt=1
while [ $attempt -le $max_attempts ]; do
    if curl -s http://localhost:8080/health > /dev/null; then
        print_status "Backend is ready!"
        break
    fi
    if [ $attempt -eq $max_attempts ]; then
        print_error "Backend failed to start within timeout"
        kill $BACKEND_PID 2>/dev/null
        exit 1
    fi
    attempt=$((attempt + 1))
    sleep 1
done

# Start frontend in background
cd web
npm run preview &
FRONTEND_PID=$!

# Function to handle cleanup on script exit
cleanup() {
    print_status "Shutting down services..."
    kill $BACKEND_PID 2>/dev/null
    kill $FRONTEND_PID 2>/dev/null
    exit 0
}

# Set up trap for cleanup
trap cleanup SIGINT SIGTERM

print_status "Services started successfully!"
print_status "Frontend: http://localhost:4173"
print_status "Backend API: http://localhost:8080/api/v1"
print_status "Press Ctrl+C to stop all services"

# Wait for user interrupt
wait 