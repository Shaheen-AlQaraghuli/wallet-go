MAKEFILE := $(abspath $(firstword $(MAKEFILE_LIST)))
MAKEFILE_DIR := $(abspath $(dir $(MAKEFILE)))
BIN_DIR := $(MAKEFILE_DIR)/bin
GOLANGCI_LINT_VERSION:=v1.64.5

include .env

install-linter:
ifeq ("$(wildcard $(BIN_DIR)/golint/$(GOLANGCI_LINT_VERSION)/golangci-lint)","")
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b "$(BIN_DIR)/golint/$(GOLANGCI_LINT_VERSION)" "$(GOLANGCI_LINT_VERSION)"
endif

install-tools: install-linter
	# Install swag CLI for Swagger doc generation
	@GOBIN="$(BIN_DIR)" go install github.com/swaggo/swag/cmd/swag@latest

	# Install mockery for mock generation
	@GOBIN="$(BIN_DIR)" go install github.com/vektra/mockery/v2@v2.20.0

	# Install goose for DB migrations
	@GOBIN="$(BIN_DIR)" go install github.com/pressly/goose/v3/cmd/goose@latest

lint: install-linter
	@"$(BIN_DIR)/golint/$(GOLANGCI_LINT_VERSION)/golangci-lint" run --fix

run:
	go run ./cmd/.

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

setup:
	docker-compose up -d

db-up:
	docker-compose up postgres -d

db-down:
	docker-compose stop postgres

migrate-up:
	@$(BIN_DIR)/goose -dir=./database/migrations postgres $(DATABASE_DSN) up

migrate-down:
	@$(BIN_DIR)/goose -dir=./database/migrations postgres $(DATABASE_DSN) down

migrate-reset:
	@$(BIN_DIR)/goose -dir=./database/migrations postgres $(DATABASE_DSN) reset

migrate-create: ## Create a new DB migration file (E.g.: make migrate-create name=create_table)
	@$(BIN_DIR)/goose -dir=./database/migrations postgres $(DATABASE_DSN) create $(name) sql

doc-gen:
	@$(BIN_DIR)/swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal --quiet

