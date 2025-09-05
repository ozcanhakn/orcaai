#!/bin/bash

# OrcaAI Setup Script
# This script helps you set up the complete OrcaAI development environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main setup function
main() {
    print_header "OrcaAI Setup Script"
    echo "This script will set up your OrcaAI development environment"
    echo ""

    # Check system requirements
    check_requirements
    
    # Create project directories
    setup_directories
    
    # Setup backend
    setup_backend
    
    # Setup Python worker
    setup_python_worker
    
    # Setup dashboard
    setup_dashboard
    
    # Setup database
    setup_database
    
    # Setup environment variables
    setup_environment
    
    # Final instructions
    print_final_instructions
}

check_requirements() {
    print_header "Checking System Requirements"
    
    # Check Go
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_status "Go found: $GO_VERSION"
    else
        print_error "Go not found. Please install Go 1.21+ from https://golang.org/dl/"
        exit 1
    fi
    
    # Check Node.js
    if command_exists node; then
        NODE_VERSION=$(node --version)
        print_status "Node.js found: $NODE_VERSION"
    else
        print_error "Node.js not found. Please install Node.js 18+ from https://nodejs.org/"
        exit 1
    fi
    
    # Check Python
    if command_exists python3; then
        PYTHON_VERSION=$(python3 --version)
        print_status "Python found: $PYTHON_VERSION"
    else
        print_error "Python3 not found. Please install Python 3.9+"
        exit 1
    fi
    
    # Check PostgreSQL
    if command_exists psql; then
        print_status "PostgreSQL client found"
    else
        print_warning "PostgreSQL client not found. You'll need to install PostgreSQL"
        echo "  - macOS: brew install postgresql"
        echo "  - Ubuntu: sudo apt-get install postgresql postgresql-contrib"
        echo "  - Windows: Download from https://www.postgresql.org/download/"
    fi
    
    # Check Redis
    if command_exists redis-cli; then
        print_status "Redis client found"
    else
        print_warning "Redis client not found. You'll need to install Redis"
        echo "  - macOS: brew install redis"
        echo "  - Ubuntu: sudo apt-get install redis-server"
        echo "  - Windows: Download from https://github.com/microsoftarchive/redis/releases"
    fi
    
    echo ""
}

setup_directories() {
    print_header "Setting up Project Directories"
    
    if [ ! -d "orcaai" ]; then
        mkdir -p orcaai
        cd orcaai
    else
        cd orcaai
    fi
    
    # Create all necessary directories
    mkdir -p backend/{config,models,orchestrator,handlers,middleware,database,utils}
    mkdir -p python-ai-worker
    mkdir -p dashboard/src/{pages,components,styles}
    mkdir -p cli/{commands,client}
    mkdir -p docs
    mkdir -p scripts
    mkdir -p tests/{backend,integration}
    mkdir -p deployment
    
    print_status "Project directories created"
}

setup_backend() {
    print_header "Setting up Go Backend"
    
    cd backend
    
    # Initialize Go module if not exists
    if [ ! -f "go.mod" ]; then
        go mod init orcaai
        print_status "Go module initialized"
    fi
    
    # Install dependencies
    print_status "Installing Go dependencies..."
    go get github.com/gin-gonic/gin@v1.9.1
    go get github.com/joho/godotenv@v1.4.0
    go get github.com/lib/pq@v1.10.9
    go get github.com/redis/go-redis/v9@v9.3.0
    go get github.com/golang-jwt/jwt/v5@v5.2.0
    go get github.com/google/uuid@v1.5.0
    go get golang.org/x/crypto@v0.17.0
    go get github.com/prometheus/client_golang@v1.17.0
    
    go mod tidy
    print_status "Go dependencies installed"
    
    cd ..
}

setup_python_worker() {
    print_header "Setting up Python AI Worker"
    
    cd python-ai-worker
    
    # Create virtual environment
    if [ ! -d "venv" ]; then
        python3 -m venv venv
        print_status "Python virtual environment created"
    fi
    
    # Activate virtual environment and install dependencies
    source venv/bin/activate
    pip install --upgrade pip
    
    # Install requirements if requirements.txt exists
    if [ -f "requirements.txt" ]; then
        pip install -r requirements.txt
        print_status "Python dependencies installed"
    else
        print_warning "requirements.txt not found. Installing basic dependencies..."
        pip install fastapi uvicorn[standard] aiohttp redis pydantic
    fi
    
    deactivate
    cd ..
}

setup_dashboard() {
    print_header "Setting up Next.js Dashboard"
    
    cd dashboard
    
    # Install npm dependencies if package.json exists
    if [ -f "package.json" ]; then
        print_status "Installing npm dependencies..."
        npm install
        print_status "Dashboard dependencies installed"
    else
        print_warning "package.json not found. Creating basic Next.js setup..."
        npm init -y
        npm install next@latest react@latest react-dom@latest
        npm install --save-dev typescript @types/react @types/node
        npm install tailwindcss autoprefixer postcss
        npm install recharts lucide-react axios date-fns
    fi
    
    cd ..
}

setup_database() {
    print_header "Setting up Database"
    
    # Check if PostgreSQL is running
    if command_exists pg_isready; then
        if pg_isready -q; then
            print_status "PostgreSQL is running"
            
            # Create database if it doesn't exist
            createdb orcaai 2>/dev/null || print_warning "Database 'orcaai' may already exist"
            
            print_status "Database setup complete"
        else
            print_warning "PostgreSQL is not running. Please start PostgreSQL:"
            echo "  - macOS: brew services start postgresql"
            echo "  - Ubuntu: sudo systemctl start postgresql"
            echo "  - Windows: Start PostgreSQL service from Services"
        fi
    else
        print_warning "PostgreSQL not found. Please install PostgreSQL first."
    fi
    
    # Check Redis
    if command_exists redis-cli; then
        if redis-cli ping > /dev/null 2>&1; then
            print_status "Redis is running"
        else
            print_warning "Redis is not running. Please start Redis:"
            echo "  - macOS: brew services start redis"
            echo "  - Ubuntu: sudo systemctl start redis"
            echo "  - Windows: Start Redis server manually"
        fi
    fi
}

setup_environment() {
    print_header "Setting up Environment Variables"
    
    if [ ! -f ".env" ]; then
        if [ -f ".env.example" ]; then
            cp .env.example .env
            print_status ".env file created from .env.example"
            print_warning "Please edit .env file with your configuration:"
            echo "  - Add your AI provider API keys"
            echo "  - Update database connection strings if needed"
            echo "  - Modify JWT secret for production"
        else
            print_warning ".env.example not found. Creating basic .env file..."
            cat > .env << EOF
# Database Configuration
DATABASE_URL=postgres://localhost/orcaai?sslmode=disable
REDIS_URL=redis://localhost:6379

# Server Configuration
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-this

# AI Provider API Keys (ADD YOUR KEYS HERE)
OPENAI_API_KEY=
CLAUDE_API_KEY=
GEMINI_API_KEY=

# Cache Configuration
CACHE_ENABLED=true
CACHE_EXPIRATION=24h

# Python Worker
PYTHON_WORKER_PORT=8001
PYTHON_WORKER_URL=http://localhost:8001

# Development
ENVIRONMENT=development
DEBUG=true
EOF
            print_status "Basic .env file created"
        fi
    else
        print_status ".env file already exists"
    fi
}

print_final_instructions() {
    print_header "Setup Complete! 🎉"
    
    echo "Your OrcaAI development environment is ready!"
    echo ""
    echo -e "${GREEN}Next steps:${NC}"
    echo "1. Edit .env file and add your AI provider API keys"
    echo "2. Make sure PostgreSQL and Redis are running"
    echo "3. Start the services:"
    echo ""
    echo -e "${BLUE}Terminal 1 - Backend:${NC}"
    echo "cd backend && go run main.go"
    echo ""
    echo -e "${BLUE}Terminal 2 - Python Worker:${NC}"
    echo "cd python-ai-worker && source venv/bin/activate && python worker.py"
    echo ""
    echo -e "${BLUE}Terminal 3 - Dashboard:${NC}"
    echo "cd dashboard && npm run dev"
    echo ""
    echo -e "${GREEN}URLs:${NC}"
    echo "- Backend API: http://localhost:8080"
    echo "- Python Worker: http://localhost:8001"
    echo "- Dashboard: http://localhost:3000"
    echo "- Health Check: http://localhost:8080/health"
    echo ""
    echo -e "${YELLOW}Important:${NC}"
    echo "- Add your OpenAI, Claude, and Gemini API keys to .env"
    echo "- The database migrations will run automatically on first start"
    echo "- Check logs for any errors during startup"
    echo ""
    echo -e "${GREEN}Happy coding! 🚀${NC}"
}

# Run main function
main "$@"