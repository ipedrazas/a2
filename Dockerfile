FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS go-builder

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

# UI builder stage
FROM node:22-alpine AS ui-builder

WORKDIR /ui

# Copy UI package files
COPY ui/package.json ui/package-lock.json* ./

RUN npm ci

# Copy UI source and build
COPY ui/ .

RUN npm run build

# Final stage - minimal image
# FROM golang:1.25-alpine

FROM alpine:3.23
# Install git (needed for cloning repos), ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy the a2 binary
COPY --from=go-builder /a2 /usr/local/bin/a2

# Copy UI assets to be served by the server
COPY --from=ui-builder /ui/dist /usr/local/share/a2/ui

# Create workspace cache directory
RUN mkdir -p /workspace/a2-cache

# Create a non-root user
RUN adduser -D -u 1000 -g 1000 a2 && \
    chown -R a2:a2 /workspace


# RUN <<EOF
# go install github.com/securego/gosec/v2/cmd/gosec@latest
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# go install golang.org/x/vuln/cmd/govulncheck@latest
# go install honnef.co/go/tools/cmd/staticcheck@2025.1.1
# go install github.com/kisielk/errcheck@latest
# go install github.com/google/go-licenses/v2@latest
# EOF

# Set working directory
WORKDIR /workspace

# Expose server port
EXPOSE 8080

# Default to check command, can be overridden with server command
ENTRYPOINT ["a2"]
CMD ["check"]

USER a2
