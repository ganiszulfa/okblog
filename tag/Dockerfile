# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 for a static binary, useful for alpine base
# -ldflags="-s -w" to strip debug information and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o /tag-service .

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /tag-service .

# Expose port (if your Fiber app listens on a specific port, e.g., 3001)
# This should match the port used in main.go and docker-compose.yml
EXPOSE 3001

# Command to run the application
CMD ["/app/tag-service"] 