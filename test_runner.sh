#!/bin/bash

# 🧪 GitHub Analytics - Test Runner
# Execute testes automatizados contra a API

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

API_BASE="http://localhost:8080/api"
HEALTH_URL="http://localhost:8080/health"

echo -e "${BLUE}🧪 GitHub Analytics - Test Runner${NC}"
echo "===================================="

# Function to check if server is running
check_server() {
    echo -e "${BLUE}🔍 Checking if server is running...${NC}"
    
    if curl -s "$HEALTH_URL" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Server is running!${NC}"
        return 0
    else
        echo -e "${RED}❌ Server is not running!${NC}"
        echo -e "${YELLOW}💡 Start the server first:${NC}"
        echo "   ./dev_server.sh"
        echo "   OR"
        echo "   make start"
        echo "   OR"
        echo "   ./start_project.sh"
        return 1
    fi
}

# Function to run API tests
run_api_tests() {
    echo ""
    echo -e "${BLUE}🔬 Running API Tests...${NC}"
    echo ""
    
    # Test 1: Health check
    echo -e "${YELLOW}1️⃣ Testing Health Check...${NC}"
    if curl -s "$HEALTH_URL" | jq . > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Health check passed${NC}"
    else
        echo -e "${RED}❌ Health check failed${NC}"
    fi
    
    # Test 2: Repositories
    echo -e "${YELLOW}2️⃣ Testing Repositories endpoint...${NC}"
    if curl -s "$API_BASE/repositories/microsoft?limit=3" | jq . > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Repositories endpoint working${NC}"
    else
        echo -e "${RED}❌ Repositories endpoint failed${NC}"
    fi
    
    # Test 3: Commit Metrics
    echo -e "${YELLOW}3️⃣ Testing Commit Metrics...${NC}"
    if curl -s "$API_BASE/metrics/commits/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21" | jq . > test_commits.json 2>&1; then
        echo -e "${GREEN}✅ Commit metrics working${NC}"
        echo -e "${BLUE}   📊 Saved to: test_commits.json${NC}"
    else
        echo -e "${RED}❌ Commit metrics failed${NC}"
    fi
    
    # Test 4: PR Metrics
    echo -e "${YELLOW}4️⃣ Testing PR Metrics...${NC}"
    if curl -s "$API_BASE/metrics/prs/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21" | jq . > test_prs.json 2>&1; then
        echo -e "${GREEN}✅ PR metrics working${NC}"
        echo -e "${BLUE}   📊 Saved to: test_prs.json${NC}"
    else
        echo -e "${RED}❌ PR metrics failed${NC}"
    fi
    
    # Test 5: Chart endpoints
    echo -e "${YELLOW}5️⃣ Testing Chart endpoints...${NC}"
    if curl -s "$API_BASE/charts/commits-timeline/microsoft/vscode?start_date=2025-07-01&end_date=2025-07-21" | jq . > test_chart.json 2>&1; then
        echo -e "${GREEN}✅ Chart endpoints working${NC}"
        echo -e "${BLUE}   📊 Saved to: test_chart.json${NC}"
    else
        echo -e "${RED}❌ Chart endpoints failed${NC}"
    fi
}

# Function to run interactive tests
run_interactive_tests() {
    echo ""
    echo -e "${BLUE}🎮 Interactive Test Mode${NC}"
    echo "========================"
    echo ""
    echo "Choose what to test:"
    echo "1) All endpoints"
    echo "2) Specific repository"
    echo "3) Custom date range"
    echo "4) Performance test"
    echo "5) Exit"
    echo ""
    read -p "Enter choice (1-5): " choice
    
    case $choice in
        1)
            run_api_tests
            ;;
        2)
            echo ""
            read -p "Enter owner: " owner
            read -p "Enter repo: " repo
            echo -e "${BLUE}🔍 Testing $owner/$repo...${NC}"
            curl -s "$API_BASE/metrics/commits/$owner/$repo?start_date=2025-07-01&end_date=2025-07-21" | jq .
            ;;
        3)
            echo ""
            read -p "Enter start date (YYYY-MM-DD): " start_date
            read -p "Enter end date (YYYY-MM-DD): " end_date
            read -p "Enter owner: " owner
            read -p "Enter repo: " repo
            echo -e "${BLUE}🔍 Testing $owner/$repo from $start_date to $end_date...${NC}"
            curl -s "$API_BASE/metrics/commits/$owner/$repo?start_date=$start_date&end_date=$end_date" | jq .
            ;;
        4)
            echo -e "${BLUE}⚡ Performance test - making 10 rapid requests...${NC}"
            for i in {1..10}; do
                echo -n "."
                curl -s "$HEALTH_URL" > /dev/null
            done
            echo ""
            echo -e "${GREEN}✅ Performance test completed${NC}"
            ;;
        5)
            echo -e "${BLUE}👋 Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Invalid choice${NC}"
            run_interactive_tests
            ;;
    esac
}

# Function to monitor logs (if Docker)
monitor_logs() {
    echo -e "${BLUE}📋 Monitoring application logs...${NC}"
    echo -e "${YELLOW}💡 Press Ctrl+C to stop monitoring${NC}"
    echo ""
    
    if command -v docker-compose &> /dev/null; then
        docker-compose logs -f
    else
        echo -e "${YELLOW}⚠️  Docker Compose not available${NC}"
        echo "Monitor the terminal where you started the server manually"
    fi
}

# Main menu
show_menu() {
    echo ""
    echo -e "${BLUE}🧪 Test Options:${NC}"
    echo "1) 🚀 Run all API tests"
    echo "2) 🎮 Interactive testing"
    echo "3) 📋 Monitor logs"
    echo "4) 🔄 Check server status"
    echo "5) 📁 Open test files folder"
    echo "6) ❌ Exit"
    echo ""
    read -p "Enter choice (1-6): " choice
    
    case $choice in
        1)
            run_api_tests
            show_menu
            ;;
        2)
            run_interactive_tests
            show_menu
            ;;
        3)
            monitor_logs
            show_menu
            ;;
        4)
            check_server
            show_menu
            ;;
        5)
            echo -e "${BLUE}📁 Generated test files:${NC}"
            ls -la test_*.json 2>/dev/null || echo "No test files found yet"
            show_menu
            ;;
        6)
            echo -e "${BLUE}👋 Goodbye!${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Invalid choice${NC}"
            show_menu
            ;;
    esac
}

# Main execution
main() {
    # Check if required tools are available
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}❌ curl is required but not installed${NC}"
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}⚠️  jq not found - JSON output will be raw${NC}"
    fi
    
    # Check server status
    if ! check_server; then
        exit 1
    fi
    
    # Show main menu
    show_menu
}

# Run main function
main
