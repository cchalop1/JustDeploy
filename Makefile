GOBIN = ./bin

all: build-web copy-web-build build-go

copy-web-build:
	cp -R web/dist internal/web/dist/
	
build-go:
	@go build -o $(GOBIN)/justdeploy ./cmd/just-deploy/main.go

build-web:
	@cd web && pnpm install && pnpm run build

run:
	@$(GOBIN)/justdeploy

clean:
	@rm -f $(GOBIN)/justdeploy
	@cd web && rm -rf dist
	@rm -r internal/web/dist
