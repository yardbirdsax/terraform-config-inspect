.PHONY: vet
vet:
	@go vet ./...

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: test-fmt
test-fmt:
	@test -z $(shell make fmt)

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: setup
setup:
	@brew bundle

.PHONY: test
test: vet test-fmt lint
	@go test ./...
