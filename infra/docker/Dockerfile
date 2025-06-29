# Stage 1: Build stage - uses full Go image with build tools.
FROM golang:tip-alpine3.22 AS builder

RUN adduser -D -g '' appuser

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/server/main.go

# Stage 2: Final stage - minimal runtime image.
FROM scratch AS deployment

COPY --from=builder /build/main /app/main

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/main"]