-include local/Makefile

.PHONY: all deps-go deps-js deps build-go build-server build-cli build-js build build-docker-dev build-docker-full lint-go revive golangci-lint test-go test-js test run run-frontend clean devenv devenv-down revive-strict protobuf help

GO = GO111MODULE=on go
GO_FILES ?= ./backend/...
VERSION ?= latest
#SH_FILES ?= $(shell find ./scripts -name *.sh)

all: build

##@ Building

build-go: ## Build all Go binaries.
	@echo "build go files"
	$(GO) run build.go build

build-debug-go: ## Build all Go binaries.
	@echo "build debug go files"
	$(GO) run build.go debug 

build: build-go 

debug: build-debug-go

run: ## Run Server
	@GO111MODULE=on ./bin/linux-amd64/godbledger-web

##@ Docker

# convenience target which looks like the other top-level build-* targets
build-docker: docker-build

docker-build:
	docker build -t godbledger-web:$(VERSION) -t godbledger-web:latest -f ./utils/Dockerfile.build-web .

docker-login:
	@$(if $(strip $(shell docker ps | grep godbledger-web)), @docker exec -it godbledger-web /bin/ash || 0, @docker run -it --rm --entrypoint /bin/ash godbledger-web:$(VERSION) )

docker-start:
	GDBL_DATA_DIR=$(GDBL_DATA_DIR) GDBL_LOG_LEVEL=$(GDBL_LOG_LEVEL) GDBL_VERSION=$(VERSION) docker-compose up

docker-stop:
	docker-compose down

docker-status:
	@$(if $(strip $(shell docker ps | grep godbledger-web)), @echo "godbledger-web is running on localhost:80", @echo "godbledger-web is not running")

docker-inspect:
	docker inspect godbledger-server

docker-logs:
	@docker logs godbledger-server

docker-logs-follow:
	@docker logs -f godbledger-server
