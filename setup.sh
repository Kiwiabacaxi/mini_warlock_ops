#!/bin/bash

echo "🚀 GitHub Analytics Backend - Setup Script"
echo "=========================================="

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📄 Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  Please edit .env file and add your GitHub token!"
    echo "   Get your token at: https://github.com/settings/tokens"
else
    echo "✅ .env file already exists"
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Build the project
echo "🔨 Building the project..."
go build -o github-analytics-backend main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🎉 Setup complete!"
    echo ""
    echo "Next steps:"
    echo "1. Edit .env file and add your GitHub token"
    echo "2. Run: ./github-analytics-backend"
    echo "3. Open: http://localhost:8080"
    echo ""
    echo "Or run directly with: go run main.go"
else
    echo "❌ Build failed. Please check for errors."
    exit 1
fi
