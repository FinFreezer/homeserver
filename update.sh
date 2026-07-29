#!/bin/bash

# Update script
sleep 2

echo "=== Starting update process ==="
echo "Killing old server..."
pkill -e -f server

sleep 2

echo "Pulling latest changes..."
git pull origin main
git reset --hard origin/main
sleep 1

echo "Starting new server..."
# Use exec to replace the script process with the server
# This keeps output in the same terminal
exec ./server.sh name password
