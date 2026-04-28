# Development environment for fmatch
# Provides a reproducible build environment with Go 1.24, make, and git.
# Usage:
#   docker build -t fmatch-dev .
#   docker run --rm -v $(pwd):/app fmatch-dev make build

FROM golang:1.24-alpine

WORKDIR /app

# Install build tools
RUN apk add --no-cache make git gcc musl-dev

# Cache dependencies layer separately for faster rebuilds
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
