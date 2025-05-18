.PHONY: mock
mock:
	docker run --rm -i -v ${PWD}:/src -w /src vektra/mockery:3.2

.PHONY: test
test:
	docker run --rm -i -v ${PWD}:/src -w /src golang:1.24.2 go test ./...

.PHONY: lint
lint:
	docker run --rm -i -v ${PWD}:/src -w /src golangci/golangci-lint:v2.0.1-alpine golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	docker run --rm -i -v ${PWD}:/src -w /src golangci/golangci-lint:v2.0.1-alpine golangci-lint run --fix ./...
