# Setup name variables for the package/tool
NAME := go-sshconfig
PKG := github.com/petems/$(NAME)

.PHONY: all
all: fmt lint test

.PHONY: fmt
fmt: ## Verifies all files have men `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | grep -v '.pb.go:' | grep -v vendor | tee /dev/stderr

.PHONY: lint
lint: ## Verifies `golangci-lint` passes
	@echo "+ $@"
	@golangci-lint run ./... 

.PHONY: cover
cover: ## Runs go test with coverage
	@for d in $(shell go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=atomic "$$d"; \
	done;

.PHONY: cover_html
cover_html: cover ## Runs go test with coverage
	@go tool cover -html=profile.out

.PHONY: test
test: ## Runs the go tests
	@echo "+ $@"
	@go test ./...