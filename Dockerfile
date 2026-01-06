FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /a2 .

# Final stage - minimal image
FROM alpine:3.19

# Install git (needed for some checks) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy the binary
COPY --from=builder /a2 /usr/local/bin/a2

# Set working directory
WORKDIR /workspace

ENTRYPOINT ["a2"]
CMD ["check"]
