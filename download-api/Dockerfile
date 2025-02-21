# Stage 1: Build the Go application
FROM golang:1.18-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o download-api ./cmd/download-api

# Stage 2: Create a minimal image with the Go binary
FROM scratch

# Copy the binary from the builder stage
COPY --from=builder /app/download-api /download-api

# Set the entrypoint to the binary
ENTRYPOINT ["/download-api"]

# Expose the port the app runs on
EXPOSE 8080
