# GitHub Analytics Backend Makefile

.PHONY: setup build run dev clean test help start stop logs docker-build docker-run docker-stop full-start

# Default target
all: setup build

# Setup project dependencies and environment
setup:
	@echo "🚀 Setting up GitHub Analytics Backend..."
	@if [ ! -f .env ]; then cp .env.example .env; echo "📄 Created .env file - please add your GitHub token!"; fi
	@go mod tidy
	@echo "✅ Setup complete!"

# Build the binary
build:
	@echo "🔨 Building binary..."
	@go build -o github-analytics-backend main.go
	@echo "✅ Build complete!"

# Run the application
run: build
	@echo "🚀 Starting GitHub Analytics Backend..."
	@./github-analytics-backend

# Run in development mode (no build)
dev:
	@echo "🔧 Running in development mode..."
	@go run main.go

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f github-analytics-backend
	@echo "✅ Clean complete!"

# Run tests (when implemented)
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Compile and run with specific GitHub token
run-with-token:
	@echo "🔑 Running with GITHUB_TOKEN environment variable..."
	@if [ -z "$(TOKEN)" ]; then echo "❌ Please provide TOKEN: make run-with-token TOKEN=your_token"; exit 1; fi
	@GITHUB_TOKEN=$(TOKEN) go run main.go

# Quick demo with docker/compose repository
demo: build
	@echo "🎬 Running demo with docker/compose repository..."
	@echo "📖 Open http://localhost:8080 in your browser"
	@echo "🔧 Try these endpoints:"
	@echo "   - GET /api/repositories/docker"
	@echo "   - GET /api/metrics/commits/docker/compose"
	@echo "   - GET /api/metrics/prs/docker/compose"
	@echo "   - GET /api/charts/contributions-leaderboard/docker/compose"
	@./github-analytics-backend

# Start the complete application (Docker Compose)
start:
	@echo "🚀 Starting GitHub Analytics with Docker Compose..."
	@if [ ! -f .env ]; then echo "❌ .env file not found! Run 'make setup' first"; exit 1; fi
	@docker-compose up -d
	@echo "✅ Application started!"
	@echo "📱 Frontend: http://localhost:8080"
	@echo "🔌 API: http://localhost:8080/api"
	@echo "📊 Health: http://localhost:8080/health"
	@echo ""
	@echo "💡 Use 'make logs' to see logs"
	@echo "🛑 Use 'make stop' to stop"

# Stop the application
stop:
	@echo "🛑 Stopping GitHub Analytics..."
	@docker-compose down
	@echo "✅ Application stopped!"

# Show application logs
logs:
	@echo "📋 Showing application logs (Ctrl+C to exit)..."
	@docker-compose logs -f

# Start the complete application locally (no Docker)
full-start: setup build
	@echo "🚀 Starting GitHub Analytics locally..."
	@if [ ! -f .env ]; then echo "❌ .env file not found! Run 'make setup' first"; exit 1; fi
	@echo "🖥️  Starting backend server..."
	@echo "📱 Frontend will be available at: http://localhost:8080"
	@echo "🔌 API will be available at: http://localhost:8080/api"
	@echo "📊 Health check: http://localhost:8080/health"
	@echo ""
	@echo "💡 Press Ctrl+C to stop the server"
	@./github-analytics-backend

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t github-analytics-backend .
	@echo "✅ Docker image built!"

# Run with Docker (single container)
docker-run: docker-build
	@echo "🐳 Running with Docker..."
	@if [ ! -f .env ]; then echo "❌ .env file not found! Run 'make setup' first"; exit 1; fi
	@docker run --rm -p 8080:8080 --env-file .env github-analytics-backend

# Stop Docker containers
docker-stop:
	@echo "🐳 Stopping Docker containers..."
	@docker-compose down
	@docker stop github-analytics-backend 2>/dev/null || true
	@echo "✅ Docker containers stopped!"

# Show help
help:
	@echo "GitHub Analytics Backend - Available commands:"
	@echo ""
	@echo "🚀 Main Commands:"
	@echo "  start              Start with Docker Compose (recommended)"
	@echo "  stop               Stop Docker Compose"
	@echo "  logs               Show application logs"
	@echo "  full-start         Start locally without Docker"
	@echo ""
	@echo "🔧 Development:"
	@echo "  setup              Setup project dependencies and .env file"
	@echo "  build              Build the binary"
	@echo "  run                Build and run the application"
	@echo "  dev                Run in development mode (go run)"
	@echo "  clean              Clean build artifacts"
	@echo "  test               Run tests"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-run         Run single Docker container"
	@echo "  docker-stop        Stop all Docker containers"
	@echo ""
	@echo "🎬 Demo:"
	@echo "  demo               Run demo with docker/compose"
	@echo "  run-with-token     Run with specific token: make run-with-token TOKEN=xxx"
	@echo ""
	@echo "🚀 Quick start:"
	@echo "  1. make setup"
	@echo "  2. Edit .env file with your GitHub token"
	@echo "  3. make start       (Docker Compose)"
	@echo "     OR"
	@echo "  3. make full-start  (Local)"

# Default help if no target specified
.DEFAULT_GOAL := help
