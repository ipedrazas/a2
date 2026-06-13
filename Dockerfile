FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS go-builder

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
FROM alpine:3.23
# Install git (needed for repo checks), ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy the a2 binary
COPY --from=go-builder /a2 /usr/local/bin/a2

# Create a non-root user
RUN adduser -D -u 1000 -g 1000 a2

# Set working directory
WORKDIR /workspace

# Run the CLI by default
ENTRYPOINT ["a2"]
CMD ["check"]

USER a2
