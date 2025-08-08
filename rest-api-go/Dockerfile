# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Final stage
FROM alpine:latest

WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /main /main

# Copy crypto material (this should be mounted as a volume in a real environment)
# COPY crypto /path/to/crypto/material

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["/main"]
