# Use the official Go image as a base
FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the Go application
COPY . .

# Build the Go application
RUN go mod tidy && go build -o server

# Expose the port
EXPOSE 8080

# Start the application
CMD ["./server"]
