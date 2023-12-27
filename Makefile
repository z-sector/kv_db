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
	go test -v -race -count=1 -coverpkg=./internal/... -coverprofile=coverage.out ./...
	go tool cover -func coverage.out