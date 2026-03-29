# Build stage for UI
FROM oven/bun:1-slim AS ui-builder

WORKDIR /app/ui

COPY ui/package.json ui/bun.lock ./
RUN bun install --frozen-lockfile

COPY ui/ ./
RUN bun run build

# Build stage for Go
FROM golang:1.26.1-alpine3.21 AS go-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built UI from previous stage
COPY --from=ui-builder /app/ui/dist ./ui/dist

# Build static binary (no CGO)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o buku .

# Final stage - scratch
FROM scratch

WORKDIR /app

# Copy binary
COPY --from=go-builder /app/buku .

# Copy SQL schema for migrations
COPY --from=go-builder /app/sql ./sql

# Create data directory (as volume mount point)
VOLUME ["/app/data"]

# Use non-root user (65534 = nobody)
USER 65534:65534

# Expose port
EXPOSE 8080

# Set environment variable (must be provided at runtime)
ENV JWT_SECRET=""

# Run the application
CMD ["./buku"]
