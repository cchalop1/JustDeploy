# Stage 1: Build frontend
FROM oven/bun:1 as frontend-builder

WORKDIR /app/web

# Copy web package files
COPY web/package.json web/bun.lockb ./

# Install frontend dependencies
RUN bun install

# Copy web source code
COPY web/ ./

# Build the frontend
RUN bun run build

# Stage 2: Build Go application
FROM golang:1.23-alpine as go-builder

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Copy built frontend from the frontend stage
COPY --from=frontend-builder /app/web/dist /app/web/dist

# Copy web build to the internal directory
RUN mkdir -p /app/internal/web/dist
RUN cp -R /app/web/dist /app/internal/web/

# Build the Go application
RUN mkdir -p bin
RUN go build -ldflags "-X cchalop1.com/deploy/internal/application.Version=$(git describe --tags --abbrev=0 || echo v0.0.0)" -o bin/justdeploy ./cmd/just-deploy/main.go

# Stage 3: Create the final image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata bash curl docker tar

# Install Nixpacks
RUN curl -sSL https://nixpacks.com/install.sh | bash

# Copy the binary from the builder stage
COPY --from=go-builder /app/bin/justdeploy /usr/local/bin/justdeploy

# Create necessary directories
RUN mkdir -p /app/data

# Expose the application port
EXPOSE 5915

# Set the entry point
ENTRYPOINT ["justdeploy"] 