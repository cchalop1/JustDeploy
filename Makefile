BINDIR 			= ./bin
EXECUTABLE  	= justdeploy
GIT_TAG 		= $(shell git describe --tags --abbrev=0)
LDFLAGS         = -ldflags "-X cchalop1.com/deploy/internal/application.Version=$(GIT_TAG)"

all: build-web copy-web-build build-go

build-ci: build-web copy-web-build build-darwin-arm build-darwin-x86 build-linux-arm build-linux-x86

build-darwin-arm:
	@env GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINDIR)/$(EXECUTABLE)-darwin-arm ./cmd/just-deploy/main.go

build-darwin-x86:
	@env GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINDIR)/$(EXECUTABLE)-darwin-x86 ./cmd/just-deploy/main.go

build-linux-arm:
	@env GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINDIR)/$(EXECUTABLE)-linux-arm ./cmd/just-deploy/main.go

build-linux-x86:
	@env GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINDIR)/$(EXECUTABLE)-linux-x86 ./cmd/just-deploy/main.go
	
build-go:
	@go build $(LDFLAGS) -o $(BINDIR)/$(EXECUTABLE) ./cmd/just-deploy/main.go

build-web:
	@cd web && bun install && bun run build

copy-web-build:
	cp -R web/dist internal/web/dist/

swagger:
	swag init --parseDependency --parseInternal -g cmd/just-deploy/main.go

run:
	@$(BINDIR)/$(EXECUTABLE)

clean:
	@rm -rf $(BINDIR)
	@cd web && rm -rf dist
	@rm -r internal/web/dist
