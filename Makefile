BIN := "./bin/app"
BIN_CLI := "./bin/cli"

.PHONY: install-lint-deps
install-lint-deps:
	@(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: install-gomock
install-gomock:
	@(which mockgen > /dev/null) || go install go.uber.org/mock/mockgen@v0.4.0

.PHONY: generate
generate: install-gomock
	go generate ./...

.PHONY: test
test:
	go test -v -race -count=1 -coverpkg=./internal/...,./config/...,./pkg/dsem/... -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

.PHONY: build
build:
	go build -v -o $(BIN) ./cmd/tcp

.PHONY: run
run: build
	$(BIN)

.PHONY: build-cli
build-cli:
	go build -v -o $(BIN_CLI) ./cmd/cli

.PHONY: run-cli
run-cli: build-cli
	$(BIN_CLI)