#!/bin/bash

# Update script
echo "Starting update..."

# Pull latest changes
git pull origin main

pkill -e server

sleep 5
# Restart server, use terminal variables.
exec ./server.sh {name} {password}
