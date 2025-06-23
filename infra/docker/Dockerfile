# Stage 1: Build stage - uses full Go image with build tools.
FROM golang:tip-alpine3.22 AS builder

# Install git and ca-certificates for dependency management and HTTPS
# Git is needed for Go modules that reference Git repositories
# ca-certificates ensure secure HTTPS connections work properly.
RUN apk add --no-cache git ca-certificates tzdata

# Create a non-root user for security.
# This user will own the application files and run the final container.
RUN adduser -D -g '' appuser

# Set the working directory inside the container.
WORKDIR /build

# Copy go.mod and go.sum first for better Docker layer caching.
COPY go.mod go.sum ./

# Download dependencies
# This step is cached unless go.mod or go.sum change.
RUN go mod download

# Verify dependencies to ensure integrity.
RUN go mod verify

# Copy the entire source code.
COPY . .

# Build the application with optimizations
# CGO_ENABLED=0: Disable CGO for a static binary (no C dependencies)
# GOOS=linux: Target Linux OS (important for cross-compilation)
# GOARCH=amd64: Target 64-bit architecture
# -ldflags="-w -s": Strip debug info and symbol table to reduce binary size
# -a: Force rebuilding of packages
# -installsuffix cgo: Add suffix to package installation directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/server/main.go

# Stage 2: Final stage - minimal runtime image.
FROM scratch AS deployment

# Import timezone data and ca-certificates from builder stage.
# scratch image has nothing, so we need to copy essential files.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the built binary from the builder stage.
COPY --from=builder /build/main /app/main

# Switch to non-root user for security.
USER appuser

# Expose port.
EXPOSE 8080

# Set the entrypoint to run the application.
ENTRYPOINT ["/app/main"]