#!/bin/bash

# Stop the current container
sudo docker compose down

# Pull latest code changes
git fetch
git pull

# Rebuild and start the container
sudo docker compose up --build -d
