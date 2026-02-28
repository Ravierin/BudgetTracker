#!/bin/bash

# BudgetTracker - Setup Script
# Creates .env with placeholder values

set -e

if [ ! -f .env ]; then
    echo "📝 Creating .env..."
    cat > .env << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=CHANGE_ME
DB_NAME=BudgetTracker
DB_SSL_MODE=disable
EOF
    echo "✅ .env created!"
    echo ""
    echo "⚠️  IMPORTANT: Edit .env and set your PostgreSQL password!"
    echo "   Then run: docker-compose up -d"
    echo ""
else
    echo "✅ .env already exists"
fi
