FROM oven/bun:1.3.10 AS frontend-builder
WORKDIR /app/web
COPY web/package.json web/bun.lock ./
RUN bun install --frozen-lockfile
COPY web/ ./
RUN rm -f .env.development && bun run build --mode production

FROM golang:1.26.10-alpine3.24 AS go-builder
WORKDIR /app
ARG VERSION=v0.0.0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/web/dist /app/internal/web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-X cchalop1.com/deploy/internal/application.Version=${VERSION}" -o /app/justdeploy ./cmd/just-deploy

FROM alpine:3.22
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata bash curl docker-cli tar git
# Nixpacks is optional. Install a verified, versioned binary on the host or supply
# a derived image; the runtime never downloads and executes an unversioned script.
COPY --from=go-builder /app/justdeploy /usr/local/bin/justdeploy
ENV JUSTDEPLOY_DATA_DIR=/app/data JUSTDEPLOY_LISTEN_ADDR=0.0.0.0
RUN mkdir -p /app/data && chmod 700 /app/data
VOLUME ["/app/data"]
EXPOSE 5915
ENTRYPOINT ["justdeploy"]
