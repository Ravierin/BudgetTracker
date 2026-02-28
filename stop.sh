#!/bin/bash

# BudgetTracker - Stop Script

set -e

echo "🛑 Stopping BudgetTracker..."
echo ""

docker-compose down

echo ""
echo "✅ All services stopped."
echo ""
echo "📝 To start again: ./start.sh"
echo ""
