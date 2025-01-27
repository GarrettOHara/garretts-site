# Use the official Go image as a base
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the Go module files first
COPY go.mod go.sum ./

# Download dependencies (this step will be cached if go.mod and go.sum haven't changed)
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the application
RUN go build -o server server.go

# Expose the port
EXPOSE 8080

# Command to start the server
CMD ["./server"]

# # Use the official Go image as a base
# FROM golang:latest
# 
# # Set the working directory
# WORKDIR /app
# 
# # Copy the application files
# COPY . .
# 
# # Build the application
# RUN go build -o server server.go
# 
# # Expose the port
# EXPOSE 8080
# 
# # Command to start the server
# CMD ["./server"]
# 
