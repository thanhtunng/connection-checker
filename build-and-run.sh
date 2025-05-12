#!/bin/bash

# Stop any running containers
echo "Stopping any running containers..."
docker-compose down

# Remove old images
echo "Removing old images..."
docker rmi ulchecker-frontend:latest ulchecker-backend:latest 2>/dev/null || true

# Build images
echo "Building images..."
docker-compose build

# Run containers
echo "Starting containers..."
docker-compose up -d

echo "Application is running!"
echo "Frontend: http://localhost:3000"
echo "Backend: http://localhost:8080" 