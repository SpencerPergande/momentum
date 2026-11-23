.DEFAULT_GOAL := check

GO_PACKAGES := ./...

.PHONY: check
check: lint test ## Lint and run unit tests (what CI runs on every pull request)

.PHONY: fmt
fmt: ## Format Go code
	gofmt -w $(shell git ls-files '*.go')

.PHONY: lint
lint: ## gofmt check, go vet and staticcheck
	@unformatted="$$(gofmt -l $$(git ls-files '*.go'))"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed:"; echo "$$unformatted"; exit 1; fi
	go vet $(GO_PACKAGES)
	go tool staticcheck $(GO_PACKAGES)

.PHONY: test
test: ## Unit and handler tests
	go test -race -count=1 $(GO_PACKAGES)

.PHONY: itest
itest: ## Integration tests against real Postgres and NATS (needs Docker)
	go test -race -count=1 -tags=integration $(GO_PACKAGES)

.PHONY: cover
cover: ## Unit tests with a coverage report
	go test -count=1 -coverprofile=coverage.out $(GO_PACKAGES)
	go tool cover -func=coverage.out | tail -1

.PHONY: help
help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-10s %s\n", $$1, $$2}'
