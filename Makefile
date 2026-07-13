SERVER_PORT ?= $(shell awk -F= '$$1 == "SERVER_PORT" { print $$2; exit }' .env 2>/dev/null)
SERVER_PORT := $(if $(SERVER_PORT),$(SERVER_PORT),4040)

.PHONY: api build env fmt help lint mocks restart run stop test

help: ## Show available commands.
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

api: ## Run API checks against a running server.
	./api.sh

build: ## Build all packages.
	go build ./...

env: ## Create .env from .env-example when it is missing.
	@test -f .env || cp .env-example .env

fmt: ## Format Go source files.
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

lint: ## Run static analysis.
	go vet ./...

mocks: ## Generate GoMock mocks from contracts.
	go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/api/public/download/contract.go -destination=internal/api/public/download/mocks/contract_mock.go -package=mocks
	go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/api/public/upload/contract.go -destination=internal/api/public/upload/mocks/contract_mock.go -package=mocks
	go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/service/uploader/contract.go -destination=internal/service/uploader/mocks/contract_mock.go -package=mocks

run: ## Start the HTTP server.
	go run ./cmd/server

stop: ## Stop the HTTP server.
	@pids="$$(lsof -tiTCP:$(SERVER_PORT) -sTCP:LISTEN)"; \
	if [ -z "$$pids" ]; then \
		echo "No server is listening on port $(SERVER_PORT)."; \
	else \
		if kill $$pids; then \
			echo "Stopped server on port $(SERVER_PORT)."; \
		else \
			echo "Could not stop server on port $(SERVER_PORT)." >&2; \
			exit 1; \
		fi; \
	fi

restart: stop run ## Restart the HTTP server.

test: ## Run all tests.
	go test ./...
