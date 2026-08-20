# ==========================================
# Stage 1: Build Stage
# ==========================================
FROM golang:1.22-alpine AS builder

# Install build-essential dependencies & certificates
RUN apk add --no-cache ca-certificates git tzdata

# Create non-root app user
RUN adduser -D -g "" -u 10001 appuser

WORKDIR /app

# Cache Go modules
COPY go.mod ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build flags for minimal static binary with injected version
ARG VERSION=1.0.0
ARG COMMIT_SHA=unknown

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X 'main.Version=${VERSION}-${COMMIT_SHA}'" \
    -trimpath \
    -o /app/bin/server ./cmd/server

# ==========================================
# Stage 2: Final Hardened Runtime Stage
# ==========================================
FROM alpine:3.19 AS final

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Copy non-root user from builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy compiled binary
COPY --from=builder --chown=10001:10001 /app/bin/server /app/server

# Set security Context & user execution
USER 10001:10001
WORKDIR /app

# OCI Image Annotations
LABEL org.opencontainers.image.title="Enterprise Go Microservice" \
      org.opencontainers.image.description="Production-grade Go API microservice with CI/CD automation" \
      org.opencontainers.image.vendor="Enterprise DevOps Team" \
      org.opencontainers.image.authors="devops@example.com" \
      org.opencontainers.image.source="https://github.com/example/ultimate-ci-cd-pipeline"

EXPOSE 8080

# Production Health check probe
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/app/server"]
