#!/bin/bash

# Peace Application Deployment Script
# Usage: ./scripts/deploy.sh [production]

set -e

ENVIRONMENT=${1:-production}
PROJECT_DIR="/opt/peace"
COMPOSE_FILE="docker-compose.${ENVIRONMENT}.yml"

echo "🚀 Deploying Peace application to ${ENVIRONMENT}..."

# Check if environment is valid
if [[ "$ENVIRONMENT" != "production" ]]; then
    echo "❌ Error: Environment must be 'production'"
    exit 1
fi

# Check if compose file exists
if [[ ! -f "$COMPOSE_FILE" ]]; then
    echo "❌ Error: $COMPOSE_FILE not found"
    exit 1
fi

# Check if .env file exists
if [[ ! -f ".env.${ENVIRONMENT}" ]]; then
    echo "❌ Error: .env.${ENVIRONMENT} not found"
    echo "Please copy env.${ENVIRONMENT}.example to .env.${ENVIRONMENT} and configure it"
    exit 1
fi

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p traefik/letsencrypt
mkdir -p traefik/logs
mkdir -p backups

# Set proper permissions for Let's Encrypt
chmod 600 traefik/letsencrypt

# Pull latest images
echo "📥 Pulling latest images..."
docker-compose -f "$COMPOSE_FILE" --env-file ".env.${ENVIRONMENT}" pull

# Stop existing containers
echo "🛑 Stopping existing containers..."
docker-compose -f "$COMPOSE_FILE" --env-file ".env.${ENVIRONMENT}" down

# Start new containers
echo "🚀 Starting new containers..."
docker-compose -f "$COMPOSE_FILE" --env-file ".env.${ENVIRONMENT}" up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 30

# Check health
echo "🏥 Checking service health..."
docker-compose -f "$COMPOSE_FILE" ps

# Clean up old images
echo "🧹 Cleaning up old images..."
docker system prune -f

echo "✅ Deployment to ${ENVIRONMENT} completed successfully!"
echo ""
echo "🌐 Services:"
if [[ "$ENVIRONMENT" == "staging" ]]; then
    echo "   Frontend: https://staging.peace.com"
    echo "   Backend:  https://api-staging.peace.com"
    echo "   Traefik:  http://staging.peace.com:8080"
else
    echo "   Frontend: https://peace.com"
    echo "   Backend:  https://api.peace.com"
    echo "   Traefik:  https://traefik.peace.com"
fi
echo ""
echo "📊 To view logs: docker-compose -f $COMPOSE_FILE logs -f"
echo "🛑 To stop: docker-compose -f $COMPOSE_FILE down"
