-include local/Makefile

.PHONY: all deps-go deps-js deps build-go build-server build-cli build-js build build-docker-dev build-docker-full lint-go revive golangci-lint test-go test-js test run run-frontend clean devenv devenv-down revive-strict protobuf help

GO = GO111MODULE=on go
GO_FILES ?= ./pkg/...
VERSION ?= latest
#SH_FILES ?= $(shell find ./scripts -name *.sh)

all: build

##@ Building

build-go: ## Build all Go binaries.
	@echo "build go files"
	$(GO) run build.go build

build: build-go 

run: ## Run Server
	@GO111MODULE=on ./bin/linux-amd64/godbledger-web

##@ Testing

test-go: ## Run tests for backend.
	@echo "test backend"
	$(GO) test -v ./pkg/...

test: test-go

##@ Linting

scripts/go/bin/revive: scripts/go/go.mod
	@cd scripts/go; \
	$(GO) build -o ./bin/revive github.com/mgechev/revive

revive: scripts/go/bin/revive
	@echo "lint via revive"
	@scripts/go/bin/revive \
		-formatter stylish \
		-config ./scripts/go/configs/revive.toml \
		$(GO_FILES)

revive-strict: scripts/go/bin/revive
	@echo "lint via revive (strict)"
	@scripts/revive-strict scripts/go/bin/revive

scripts/go/bin/golangci-lint: scripts/go/go.mod
	@cd scripts/go; \
	$(GO) build -o ./bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

golangci-lint: scripts/go/bin/golangci-lint
	@echo "lint via golangci-lint"
	@scripts/go/bin/golangci-lint run \
		--config ./scripts/go/configs/.golangci.toml \
		$(GO_FILES)

lint-go: golangci-lint revive revive-strict # Run all code checks for backend.

##@ Docker

# convenience target which looks like the other top-level build-* targets
build-docker: docker-build

docker-build:
	docker build -t godbledger-web:$(VERSION) -t godbledger-web:latest -f ./utils/Dockerfile.build-web .

docker-login:
	@$(if $(strip $(shell docker ps | grep godbledger-server)), @docker exec -it godbledger-server /bin/ash || 0, @docker run -it --rm --entrypoint /bin/ash godbledger:$(VERSION) )

docker-start:
	GDBL_DATA_DIR=$(GDBL_DATA_DIR) GDBL_LOG_LEVEL=$(GDBL_LOG_LEVEL) GDBL_VERSION=$(VERSION) docker-compose up

docker-stop:
	docker-compose down

docker-status:
	@$(if $(strip $(shell docker ps | grep godbledger-server)), @echo "godbledger-server is running on localhost:50051", @echo "godbledger-server is not running")

docker-inspect:
	docker inspect godbledger-server

docker-logs:
	@docker logs godbledger-server

docker-logs-follow:
	@docker logs -f godbledger-server

##@ Helpers

clean: ## Clean up intermediate build artifacts.
	@echo "cleaning"
	rm -rf node_modules
	rm -rf public/build

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
