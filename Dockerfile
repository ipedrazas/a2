FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

# Build arguments with defaults
ARG VERSION=dev
ARG GITSHA=unknown
ARG BUILDDATE=unknown

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files first for caching
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build the binary for the target platform
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -X github.com/ipedrazas/a2/pkg/version.Version=${VERSION} -X github.com/ipedrazas/a2/pkg/version.GitSHA=${GITSHA} -X github.com/ipedrazas/a2/pkg/version.BuildDate=${BUILDDATE}" \
    -o /a2 .

# Final stage - minimal image
FROM golang:1.25-alpine

# Install git (needed for some checks) and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy the binary
COPY --from=builder /a2 /usr/local/bin/a2

# Install tools
RUN <<EOF
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install honnef.co/go/tools/cmd/staticcheck@2025.1.1
go install github.com/kisielk/errcheck@latest
# Add other tools here
EOF

# Set working directory
WORKDIR /workspace

# Create a non-root user
RUN adduser -D -u 1000 -g 1000 a2

ENTRYPOINT ["a2"]
CMD ["check"]

USER a2
