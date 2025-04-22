# FROM golang:latest
# 
# WORKDIR /app
# 
# # Copy only go.mod and go.sum first to cache dependencies
# COPY go.mod go.sum ./
# 
# # Download dependencies (cache this step)
# RUN go mod download
# 
# # Copy the rest of the application
# COPY . .
# 
# # Use a minimal base image for the final container
# FROM alpine:latest  
# 
# WORKDIR /root/
# 
# # Copy the built binary from the builder stage
# COPY --from=builder /app/server .
# 
# # Run the binary
# CMD ["./server"]

# # First stage: build
# FROM golang:latest AS builder
# 
# WORKDIR /app
# 
# # Copy go mod files and download dependencies
# COPY go.mod go.sum ./
# RUN go mod download
# 
# # Copy rest of the source code
# COPY . .
# 
# # Build the Go application
# RUN go build -o server .
# 
# # Second stage: run
# FROM alpine:latest
# 
# WORKDIR /root/
# 
# # Copy built binary from the builder stage
# COPY --from=builder /app/server .
# 
# CMD ["./server"]
# # First stage: build
# FROM golang:latest AS builder
# 
# WORKDIR /app
# 
# COPY go.mod go.sum ./
# RUN go mod download
# 
# COPY . .
# 
# # 🔥 Build from the correct subdirectory
# RUN go build -o server ./cmd/server
# 
# # Second stage: run
# FROM alpine:latest
# 
# WORKDIR /root/
# 
# COPY --from=builder /app/server .
# 
# CMD ["./server"]


# # Build stage
# FROM golang:latest AS builder
# 
# WORKDIR /app
# 
# COPY go.mod go.sum ./
# RUN go mod download
# 
# COPY . .
# 
# # ✅ Cross-compile for Linux
# RUN GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
# 
# # Runtime stage
# FROM alpine:latest
# 
# WORKDIR /root/
# 
# # Copy the binary from builder
# COPY --from=builder /app/server .
# 
# # ✅ Ensure executable
# RUN chmod +x ./server
# 
# CMD ["./server"]

# # Build stage
# FROM golang:latest AS builder
# 
# WORKDIR /app
# 
# COPY go.mod go.sum ./
# RUN go mod download
# 
# COPY . .
# 
# # ✅ Static binary for Alpine
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
# 
# # Runtime stage
# FROM alpine:latest
# 
# WORKDIR /root/
# 
# COPY --from=builder /app/server .
# 
# # Optional: ensure it's executable
# RUN chmod +x ./server
# 
# CMD ["./server"]

FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
