#!/bin/bash

# 🚀 GitHub Analytics - Complete Project Startup Script
# Este script configura e executa o projeto completo

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}🚀 GitHub Analytics Backend${NC}"
    echo -e "${BLUE}================================${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Check dependencies
check_dependencies() {
    print_info "Checking dependencies..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+ first."
        echo "Visit: https://golang.org/doc/install"
        exit 1
    fi
    print_success "Go is installed: $(go version)"
    
    # Check Docker (optional)
    if command -v docker &> /dev/null; then
        print_success "Docker is available: $(docker --version | head -1)"
        DOCKER_AVAILABLE=true
    else
        print_warning "Docker not found - will run locally only"
        DOCKER_AVAILABLE=false
    fi
    
    # Check Docker Compose (optional)
    if command -v docker-compose &> /dev/null; then
        print_success "Docker Compose is available: $(docker-compose --version)"
        COMPOSE_AVAILABLE=true
    else
        print_warning "Docker Compose not found - will run locally only"
        COMPOSE_AVAILABLE=false
    fi
}

# Setup environment
setup_environment() {
    print_info "Setting up environment..."
    
    # Create .env if not exists
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            print_success "Created .env file from template"
        else
            cat > .env << EOF
GITHUB_TOKEN=your_github_token_here
DATE_RANGE_START=2025-07-01
DATE_RANGE_END=2025-07-31
PORT=8080
EOF
            print_success "Created .env file with defaults"
        fi
        print_warning "Please edit .env file and add your GitHub token!"
        print_info "Get your token at: https://github.com/settings/tokens"
        echo ""
        read -p "Press Enter after you've added your GitHub token to .env file..."
    else
        print_success ".env file already exists"
    fi
    
    # Install Go dependencies
    print_info "Installing Go dependencies..."
    go mod tidy
    print_success "Dependencies installed"
}

# Validate configuration
validate_config() {
    print_info "Validating configuration..."
    
    if [ -f .env ]; then
        source .env
        if [ -z "$GITHUB_TOKEN" ] || [ "$GITHUB_TOKEN" = "your_github_token_here" ]; then
            print_error "GITHUB_TOKEN not set in .env file!"
            print_info "Please edit .env file and add your GitHub token"
            print_info "Get your token at: https://github.com/settings/tokens"
            exit 1
        fi
        print_success "GitHub token configured"
    else
        print_error ".env file not found!"
        exit 1
    fi
}

# Build project
build_project() {
    print_info "Building project..."
    go build -o github-analytics-backend main.go
    if [ $? -eq 0 ]; then
        print_success "Build completed successfully"
    else
        print_error "Build failed"
        exit 1
    fi
}

# Start with Docker Compose
start_with_docker_compose() {
    print_info "Starting with Docker Compose..."
    
    docker-compose up -d
    
    # Wait for service to be ready
    print_info "Waiting for service to be ready..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            print_success "Service is ready!"
            break
        fi
        sleep 2
        echo -n "."
    done
    echo ""
    
    print_success "🎉 Application started successfully!"
    echo ""
    echo "📱 Frontend: http://localhost:8080"
    echo "🔌 API: http://localhost:8080/api"
    echo "📊 Health: http://localhost:8080/health"
    echo ""
    echo "Commands:"
    echo "  📋 View logs: docker-compose logs -f"
    echo "  🛑 Stop: docker-compose down"
    echo "  🔄 Restart: docker-compose restart"
}

# Start locally
start_locally() {
    print_info "Starting locally..."
    
    print_success "🎉 Starting GitHub Analytics Backend..."
    echo ""
    echo "📱 Frontend: http://localhost:8080"
    echo "🔌 API: http://localhost:8080/api"
    echo "📊 Health: http://localhost:8080/health"
    echo ""
    echo "💡 Press Ctrl+C to stop the server"
    echo ""
    
    # Run the application
    ./github-analytics-backend
}

# Main menu
show_menu() {
    echo ""
    print_info "Choose how to run the application:"
    echo ""
    echo "1) 🐳 Docker Compose (recommended)"
    echo "2) 🖥️  Local binary"
    echo "3) 🔧 Development mode (go run)"
    echo "4) 🛠️  Setup only"
    echo "5) ❌ Exit"
    echo ""
    read -p "Enter your choice (1-5): " choice
}

# Handle user choice
handle_choice() {
    case $choice in
        1)
            if [ "$DOCKER_AVAILABLE" = true ] && [ "$COMPOSE_AVAILABLE" = true ]; then
                start_with_docker_compose
            else
                print_error "Docker or Docker Compose not available!"
                print_info "Falling back to local execution..."
                build_project
                start_locally
            fi
            ;;
        2)
            build_project
            start_locally
            ;;
        3)
            print_info "Starting in development mode..."
            print_success "🎉 Starting GitHub Analytics Backend (Dev Mode)..."
            echo ""
            echo "📱 Frontend: http://localhost:8080"
            echo "🔌 API: http://localhost:8080/api"
            echo "📊 Health: http://localhost:8080/health"
            echo ""
            echo "💡 Press Ctrl+C to stop the server"
            echo ""
            go run main.go
            ;;
        4)
            print_success "Setup completed! Use './start_project.sh' to run again."
            ;;
        5)
            print_info "Goodbye! 👋"
            exit 0
            ;;
        *)
            print_error "Invalid choice. Please enter 1-5."
            show_menu
            handle_choice
            ;;
    esac
}

# Main execution
main() {
    print_header
    check_dependencies
    setup_environment
    validate_config
    
    show_menu
    handle_choice
}

# Run main function
main
