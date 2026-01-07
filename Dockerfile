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

# Install git (needed for some checks), ca-certificates, wget, and Python for some tools
RUN apk add --no-cache git ca-certificates 

# Copy the binary
COPY --from=builder /a2 /usr/local/bin/a2

# Install Go tools
RUN <<EOF
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/vuln/cmd/govulncheck@latest
go install honnef.co/go/tools/cmd/staticcheck@2025.1.1
go install github.com/kisielk/errcheck@latest
go install github.com/zricethezav/gitleaks/v8@latest
EOF

# Install Python tools
# RUN pip3 install --no-cache-dir --break-system-packages bandit semgrep

# Create a non-root user and set up Go directories
RUN adduser -D -u 1000 -g 1000 a2 && \
    mkdir -p /home/a2/go /home/a2/.cache && \
    chown -R a2:a2 /home/a2

# Set working directory
WORKDIR /workspace

# Environment for non-root Go usage
ENV GOPATH=/home/a2/go
ENV GOCACHE=/home/a2/.cache/go-build
ENV GOMODCACHE=/home/a2/go/pkg/mod
ENV PATH="/home/a2/go/bin:/go/bin:${PATH}"

ENTRYPOINT ["a2"]
CMD ["check"]

USER a2
