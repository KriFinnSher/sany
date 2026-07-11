.PHONY: api build fmt help lint mocks run test

help: ## Show available commands.
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-10s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

api: ## Run API checks against a running server.
	./api.sh

build: ## Build all packages.
	go build ./...

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

test: ## Run all tests.
	go test ./...
