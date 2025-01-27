# Use the official Go image as a base
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the application files
COPY . .

# Build the application
RUN go build -o server server.go

# Expose the port
EXPOSE 8080

# Command to start the server
CMD ["./server"]

