#!/bin/bash

# 🔧 GitHub Analytics - Development Mode
# Este script é otimizado para desenvolvimento com hot reload

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}🔧 GitHub Analytics - Development Mode${NC}"
echo "========================================"

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found! Creating from example...${NC}"
    if [ -f .env.example ]; then
        cp .env.example .env
    else
        cat > .env << EOF
GITHUB_TOKEN=your_github_token_here
DATE_RANGE_START=2025-07-01
DATE_RANGE_END=2025-07-31
PORT=8080
EOF
    fi
    echo -e "${YELLOW}📝 Please edit .env and add your GitHub token!${NC}"
    echo "Get token at: https://github.com/settings/tokens"
    echo ""
    read -p "Press Enter after adding your token..."
fi

# Install dependencies if needed
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo -e "${BLUE}📦 Installing dependencies...${NC}"
    go mod tidy
fi

echo ""
echo -e "${GREEN}🚀 Starting backend server...${NC}"
echo "📱 Frontend: http://localhost:8080"
echo "🔌 API: http://localhost:8080/api"
echo "📊 Health: http://localhost:8080/health"
echo ""
echo -e "${YELLOW}💡 This runs with live reload - modify files and restart manually${NC}"
echo -e "${YELLOW}🛑 Press Ctrl+C to stop${NC}"
echo ""

# Run with hot reload using go run
while true; do
    go run main.go
    echo ""
    echo -e "${YELLOW}🔄 Server stopped. Press Enter to restart or Ctrl+C to exit...${NC}"
    read -r
done
