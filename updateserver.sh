#!/bin/bash

# Update script
echo "Starting update..."

# Pull latest changes
git pull origin main

# Restart server, use terminal variables.
./server.sh $1 $2
