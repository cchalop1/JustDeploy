FROM oven/bun:latest as frontend-builder

WORKDIR /app/web

COPY web/package.json web/bun.lockb ./

RUN bun install

COPY web/ ./

RUN bun run build --mode production

FROM golang:1.23-alpine as go-builder

WORKDIR /app

# Add build argument for version
ARG VERSION=v0.0.0

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY --from=frontend-builder /app/web/dist /app/web/dist

RUN mkdir -p /app/internal/web/dist
RUN cp -R /app/web/dist /app/internal/web/

RUN mkdir -p bin
RUN go build -ldflags "-X cchalop1.com/deploy/internal/application.Version=${VERSION}" -o bin/justdeploy ./cmd/just-deploy/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata bash curl docker tar

RUN curl -sSL https://nixpacks.com/install.sh | bash

COPY --from=go-builder /app/bin/justdeploy /usr/local/bin/justdeploy

RUN mkdir -p /app/data

EXPOSE 5915

ENTRYPOINT ["justdeploy"] 