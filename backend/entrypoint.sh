#!/bin/sh
# Docker entrypoint script that applies migrations then starts the app

set -e

echo "🔧 Applying database migrations..."

# Wait for database to be ready
until pg_isready -h postgres -U postgres > /dev/null 2>&1; do
    echo "   Waiting for database..."
    sleep 1
done

# Apply migrations in order
export PGPASSWORD="${DB_PASSWORD:-postgres}"
psql -h postgres -U postgres -d BudgetTracker -v ON_ERROR_STOP=1 << 'EOSQL'
CREATE TABLE IF NOT EXISTS position (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    exchange VARCHAR(50) NOT NULL,
    symbol VARCHAR(50) NOT NULL,
    volume DECIMAL(20, 8) NOT NULL,
    leverage INTEGER NOT NULL,
    closed_pnl DECIMAL(20, 8) NOT NULL,
    side VARCHAR(20) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS withdrawal (
    id SERIAL PRIMARY KEY,
    exchange VARCHAR(50) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS monthly_income (
    id SERIAL PRIMARY KEY,
    exchange VARCHAR(50) NOT NULL,
    amount DECIMAL(20, 8) NOT NULL,
    pnl DECIMAL(20, 8) NOT NULL,
    date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    exchange VARCHAR(50) NOT NULL UNIQUE,
    api_key VARCHAR(500) NOT NULL,
    api_secret VARCHAR(500) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO api_keys (exchange, api_key, api_secret, is_active)
VALUES
    ('mexc', '', '', false),
    ('bybit', '', '', false)
ON CONFLICT (exchange) DO NOTHING;
EOSQL

echo "✅ Migrations applied!"
echo "🚀 Starting application..."

# Start the application
exec ./budget-tracker
