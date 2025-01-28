#!/bin/bash

# Stop the current container
docker compose down

# Pull latest code changes
git fetch
git pull

# Rebuild and start the container
docker compose up --build -d
