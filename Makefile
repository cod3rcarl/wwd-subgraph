GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: run runContainer runApp dbNew tests test coverage htmlcover bench gomodtidy cilint cilint-fix gofumpt check fix gen genAll help

DEFAULT: help

ifeq ($(OS),Windows_NT)
EXT = .exe
else
EXT =
endif

127.0.0.1:50060

run: # 50060
	GRPC_HOST=0.0.0.0 \
	GRPC_PORT=50050 \
	DB_CONNECTION="user=mnjghfpgupfubo password=2fcdf645897feac50ddb34cd0ae211e064636c95f8453899ed289eeb183ae70a host=ec2-54-155-112-143.eu-west-1.compute.amazonaws.com port=5432 dbname=d3leb3ukf4fo82 pool_max_conns=10" \
	go run cmd/app/main.go

runContainer: ## Run the app in Docker
	docker compose run --rm --service-ports app bash

runApp: ## Run the app in Docker
	docker compose up -d app

tests: ## Run all tests in the project
	gotest -v -cover -race ./...

test: ## Run one particular test
	gotest -v -cover -race -run=$(name)

coverage: ## Run tests and export coverage
	gotest -v -coverprofile .coverage.out ./...; go tool cover -func=.coverage.out

htmlcover: coverage ## Run tests and export coverage in html format
	go tool cover -html=.coverage.out -o=.coverage.html

bench: ## Run bench tests
	gotest -run=Bench -bench=. ./...

gomodtidy: ## Tidy up Go modules
	go mod tidy

cilint: ## Run linters
	golangci-lint$(EXT) run -v

cilint-fix: ## Run linters with auto fixes
	golangci-lint$(EXT) run --fix -v

gofumpt: ## Format code with gofumpt
	gofumpt$(EXT) -l -w .

check: gofumpt cilint gomodtidy ## Format code, run linters and tidy up Go modules

fix: gofumpt cilint-fix gomodtidy ## Format code, run linters with auto fixes and tidy up Go modules

gen: ## Generate GraphQL code
	go get -d github.com/99designs/gqlgen; go run github.com/99designs/gqlgen generate

genAll: ## Generate everything in the project
	go generate ./...

help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
