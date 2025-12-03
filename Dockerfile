# ---- Build Stage ----
FROM golang:1.24 AS builder

WORKDIR /app

# Copy mod files
COPY go.mod go.sum ./

# Clean and download deps
RUN go mod tidy
RUN go mod download

# Copy everything else
COPY . .

# Build binary
RUN go build -o server ./cmd/server

# ---- Run Stage ----
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
