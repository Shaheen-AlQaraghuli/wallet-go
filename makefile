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
	# Install swaggo for API doc generation
	@GOBIN="$(BIN_DIR)" go install github.com/swaggo/swag/cmd/swag@latest
	# Install redocly CLI for API doc manipulation
	@rm -rf /opt/homebrew/lib/node_modules/@redocly
	@npm install -g @redocly/cli || sudo npm install -g @redocly/cli
	# Install speccy for linting API doc
	@rm -rf /opt/homebrew/lib/node_modules/speccy
	@npm install -g speccy || sudo npm install -g speccy
	# Install mockery for mock generation
	@GOBIN="$(BIN_DIR)" go install github.com/vektra/mockery/v2@v2.20.0
	# Install goose for DB migrations
	@GOBIN="$(BIN_DIR)" go install github.com/pressly/goose/v3/cmd/goose@latest
	# Install gci to automatically format imports
	@GOBIN="$(BIN_DIR)" go install github.com/daixiang0/gci@latest

lint: install-linter
	@"$(BIN_DIR)/golint/$(GOLANGCI_LINT_VERSION)/golangci-lint" run --fix

run:
	go run ./cmd/.

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
	# Remove all temp files that might be there because of a previously failed doc-gen.
	@rm -rf ./docs/tmp
	# Generate OpenAPI v2 doc from swaggo/swag annotations.
	@$(BIN_DIR)/swag init -g ./cmd/main.go -o ./docs/tmp --parseDependency --parseInternal --quiet --collectionFormat multi
	# Convert the generated OpenAPI v2 yaml file to OpenAPI v3 yaml file.
	@docker run --rm -u $(shell id -u):$(shell id -g) -v $(PWD)/docs/tmp:/work openapitools/openapi-generator-cli:latest-release \
        generate -i /work/swagger.yaml -o /work/v3 -g openapi-yaml --minimal-update 1> /dev/null
	# Remove the path prefix from the generated schema names.
	@sed -i -e "s/gitlab_intelligentb_com_.*\.//g" ./docs/tmp/v3/openapi/openapi.yaml
	@sed -i -e "s/gitlab_intelligentb_com_.*_models_//g" ./docs/tmp/v3/openapi/openapi.yaml
	@for prefix in types apierror model models pagination time; do \
		sed -i -e "s/$$prefix\.//g" ./docs/tmp/v3/openapi/openapi.yaml ; \
    done
	@sed -i -e "s/type\: object//g" ./docs/tmp/v3/openapi/openapi.yaml
	@sleep 1
	# Replace the servers section of the above temp file. This is because swaggo/swag only supports OpenAPI v2 and v2 doesn't support multiple servers.
	@docker run --security-opt=no-new-privileges --cap-drop all --network none --rm -v $(PWD)/docs:/work mikefarah/yq '. *n load("/work/tmp/v3/openapi/openapi.yaml")' /work/overriding-template.yml > $(PWD)/docs/api.yml
	@sleep 1
	# Remove all temp files.
	@rm -rf ./docs/tmp
	# Check the final API doc.
	@docker run --rm -v $(PWD)/docs:/docs wework/speccy lint -v /docs/api.yml

doc-lint:
	@docker run --rm -v $(PWD)/docs:/docs wework/speccy lint -v /docs/api.yml
