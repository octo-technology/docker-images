# Build stage
FROM docker.io/library/golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY main.go .

ARG VERSION=dev
ARG BUILD_TIMESTAMP

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildDate=${BUILD_TIMESTAMP}" \
    -o metrics-app .

# Final image
FROM scratch

LABEL org.opencontainers.image.title="Metrics App"
LABEL org.opencontainers.image.authors="Joseph Page <josephpage@users.noreply.github.com>"
LABEL org.opencontainers.image.description="Demo app for OCTO training sessions"
LABEL org.opencontainers.image.url="https://www.octo.academy"
LABEL org.opencontainers.image.vendor="OCTO Technology"
LABEL org.opencontainers.image.source="https://github.com/octo-technology/docker-images/tree/main/metrics-app"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_TIMESTAMP}"
LABEL org.opencontainers.image.base.name="docker.io/library/scratch"

COPY --from=builder /build/metrics-app /app

EXPOSE 8080

ENTRYPOINT ["/app"]
