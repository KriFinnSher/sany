.PHONY: api build fmt lint mocks run test

api:
	./api.sh

build:
	go build ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

lint:
	go vet ./...

mocks:
	go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/api/public/download/contract.go -destination=internal/api/public/download/mocks/contract_mock.go -package=mocks
	go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/api/public/upload/contract.go -destination=internal/api/public/upload/mocks/contract_mock.go -package=mocks
	go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/service/uploader/contract.go -destination=internal/service/uploader/mocks/contract_mock.go -package=mocks

run:
	go run ./cmd/server

test:
	go test ./...
