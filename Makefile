# GitHub Analytics Backend Makefile

.PHONY: setup build run dev clean test help

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

# Show help
help:
	@echo "GitHub Analytics Backend - Available commands:"
	@echo ""
	@echo "  setup              Setup project dependencies and .env file"
	@echo "  build              Build the binary"
	@echo "  run                Build and run the application"
	@echo "  dev                Run in development mode (go run)"
	@echo "  clean              Clean build artifacts"
	@echo "  test               Run tests"
	@echo "  demo               Run demo with docker/compose"
	@echo "  run-with-token     Run with specific token: make run-with-token TOKEN=xxx"
	@echo "  help               Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  1. make setup"
	@echo "  2. Edit .env file with your GitHub token"
	@echo "  3. make demo"

# Default help if no target specified
.DEFAULT_GOAL := help
