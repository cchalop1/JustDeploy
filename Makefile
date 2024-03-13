GOBIN = ./bin

all: build-web copy-web-build build-go

copy-web-build:
	cp -R web/dist internal/application/
	
build-go:
	@go build -o $(GOBIN)/just-deploy ./cmd/just-deploy/main.go

build-web:
	@cd web && pnpm install && pnpm run build

run:
	@$(GOBIN)/just-deploy

clean:
	@rm -f $(GOBIN)/just-deploy
	@cd web && rm -rf dist
