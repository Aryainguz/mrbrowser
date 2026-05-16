# ─────────────────────────────────────────────────────────────
# Stage 1: Build Go binary
# ─────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Install CGO dependencies (required for go-sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /build

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /build/mr-browser ./cli

# ─────────────────────────────────────────────────────────────
# Stage 2: Runtime image
# ─────────────────────────────────────────────────────────────
FROM debian:bookworm-slim AS runtime

# Install Chromium and required libraries
RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium \
    chromium-sandbox \
    ca-certificates \
    fonts-liberation \
    fonts-noto-color-emoji \
    libnss3 \
    libatk-bridge2.0-0 \
    libgtk-3-0 \
    libxss1 \
    libasound2 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r mrbrowser && useradd -r -g mrbrowser -d /app -s /sbin/nologin mrbrowser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/mr-browser /usr/local/bin/mr-browser

# Create data directory
RUN mkdir -p /app/data /app/screenshots && chown -R mrbrowser:mrbrowser /app

# Switch to non-root user
USER mrbrowser

# Environment defaults
ENV MRBROWSER_PORT=7331 \
    MRBROWSER_HEADLESS=true \
    MRBROWSER_NO_SANDBOX=true \
    MRBROWSER_CHROMIUM_PATH=/usr/bin/chromium \
    MRBROWSER_DB_PATH=/app/data/mrbrowser.db \
    MRBROWSER_LOG_LEVEL=info

EXPOSE 7331

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD mr-browser --help > /dev/null || exit 1

ENTRYPOINT ["mr-browser"]
CMD ["--help"]
