#!/bin/bash

# Script to setup local host entries for development
# Run with: sudo bash scripts/setup-hosts.sh

HOSTS_FILE="/etc/hosts"
BACKUP_FILE="/etc/hosts.backup.$(date +%Y%m%d_%H%M%S)"

# Create backup
echo "Creating backup of hosts file..."
cp "$HOSTS_FILE" "$BACKUP_FILE"

# Define host entries
ENTRIES=(
    "127.0.0.1 localhost"
    "127.0.0.1 api.localhost"
    "127.0.0.1 ws.localhost"
    "127.0.0.1 traefik.localhost"
)

echo "Adding host entries for Peace development environment..."

for entry in "${ENTRIES[@]}"; do
    if ! grep -q "$entry" "$HOSTS_FILE"; then
        echo "Adding: $entry"
        echo "$entry" >> "$HOSTS_FILE"
    else
        echo "Already exists: $entry"
    fi
done

echo ""
echo "Setup complete! Host entries added:"
echo "- http://localhost (Frontend)"
echo "- http://api.localhost/api (Backend API)"
echo "- ws://ws.localhost/ws (WebSocket)"
echo "- http://traefik.localhost (Traefik Dashboard, also available at http://localhost:8090)"
echo ""
echo "You can now start the services with: make dev"
