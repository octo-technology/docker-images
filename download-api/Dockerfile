# Stage 1: Build the Go application
FROM golang:1.26-alpine AS builder

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

# Labels OCI pour identification de l'image
ARG VERSION="dev"
ARG BUILD_TIMESTAMP="unknown"

LABEL org.opencontainers.image.title="Download API"
LABEL org.opencontainers.image.authors="Joseph Page <josephpage@users.noreply.github.com>, Yannig Perre <yannig@users.noreply.github.com>"
LABEL org.opencontainers.image.description="API for downloading files, used in OCTO training sessions"
LABEL org.opencontainers.image.url="https://www.octo.academy"
LABEL org.opencontainers.image.vendor="OCTO Technology"
LABEL org.opencontainers.image.source="https://github.com/octo-technology/docker-images/tree/main/download-api"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_TIMESTAMP}"
LABEL org.opencontainers.image.base.name="docker.io/library/scratch"

# Copy the binary from the builder stage
COPY --from=builder /app/download-api /download-api

# Set the entrypoint to the binary
ENTRYPOINT ["/download-api"]

# Expose the port the app runs on
EXPOSE 8080
