#!/bin/bash

# Update script
echo "Starting update..."

sleep 5

# Pull latest changes
git pull origin main

sleep 5
# Restart server, use terminal variables.
exec ./server.sh {name} {password}
