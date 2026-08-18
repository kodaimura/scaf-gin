DOCKER_COMPOSE := docker compose
ENV ?= dev
PROJECT_NAME ?= $(notdir $(CURDIR))
DOCKER_COMPOSE_FILE := $(if $(filter prod,$(ENV)),-f docker-compose.prod.yml,-f docker-compose.yml)
DOCKER_COMPOSE_CMD := $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE)
E2E_COMPOSE_CMD := $(DOCKER_COMPOSE) -p scaf-gin-e2e -f docker-compose.yml -f docker-compose.test.yml
API_SERVICE := api
MIGRATE_SERVICE := migrate

.DEFAULT_GOAL := help

.PHONY: init up build build_no_cache build_prod down down_volumes stop exec shell logs ps reup check lint format format_check test test_e2e smoke routes migrate current history help

## -----------------------------
## Base Commands
## -----------------------------

init:
	./bin/scaf-init "$(PROJECT_NAME)"

up:
	$(DOCKER_COMPOSE_CMD) up -d

build:
	$(DOCKER_COMPOSE_CMD) build

build_no_cache:
	$(DOCKER_COMPOSE_CMD) build --no-cache

build_prod:
	docker build --target runtime --tag "$(PROJECT_NAME)-runtime" .

down:
	$(DOCKER_COMPOSE_CMD) down

down_volumes:
	$(DOCKER_COMPOSE_CMD) down -v

stop:
	$(DOCKER_COMPOSE_CMD) stop

exec:
	$(DOCKER_COMPOSE_CMD) exec $(API_SERVICE) /bin/sh

shell:
	$(DOCKER_COMPOSE_CMD) run --rm $(API_SERVICE) /bin/sh

logs:
	$(DOCKER_COMPOSE_CMD) logs -f $(API_SERVICE)

ps:
	$(DOCKER_COMPOSE_CMD) ps

reup: down up

check: format_check lint test

lint:
	$(DOCKER_COMPOSE_CMD) run --rm --no-deps $(API_SERVICE) go vet ./...

format:
	$(DOCKER_COMPOSE_CMD) run --rm --no-deps $(API_SERVICE) sh -c 'gofmt -w $$(find . -type f -name "*.go")'

format_check:
	$(DOCKER_COMPOSE_CMD) run --rm --no-deps $(API_SERVICE) sh -c 'files="$$(gofmt -l $$(find . -type f -name "*.go"))"; test -z "$$files" || { echo "$$files"; exit 1; }'

test:
	$(DOCKER_COMPOSE_CMD) run --rm --no-deps $(API_SERVICE) go test ./...

test_e2e:
	@set -eu; \
	cleanup() { $(E2E_COMPOSE_CMD) down -v --remove-orphans >/dev/null 2>&1 || true; }; \
	trap cleanup EXIT INT TERM; \
	cleanup; \
	$(E2E_COMPOSE_CMD) --profile tools run --rm --build $(MIGRATE_SERVICE); \
	if ! $(E2E_COMPOSE_CMD) --profile test run --rm --build api-test; then \
		$(E2E_COMPOSE_CMD) logs --no-color $(API_SERVICE) db mailhog >&2 || true; \
		exit 1; \
	fi

smoke:
	$(DOCKER_COMPOSE_CMD) exec -T $(API_SERVICE) sh -c 'if command -v go >/dev/null 2>&1; then go run ./cmd/healthcheck; else /app/healthcheck; fi'

routes:
	$(DOCKER_COMPOSE_CMD) run --rm -e PRINT_ROUTES=true $(API_SERVICE) sh -c 'if command -v go >/dev/null 2>&1; then go run ./cmd; else /app/api; fi'

## -----------------------------
## Migrations
## -----------------------------

migrate:
	$(DOCKER_COMPOSE_CMD) run --rm $(MIGRATE_SERVICE)

current:
	$(DOCKER_COMPOSE_CMD) run --rm $(MIGRATE_SERVICE) sh -c 'if command -v go >/dev/null 2>&1; then go run ./cmd/migrate current; else /app/migrate current; fi'

history:
	$(DOCKER_COMPOSE_CMD) run --rm $(MIGRATE_SERVICE) sh -c 'if command -v go >/dev/null 2>&1; then go run ./cmd/migrate history; else /app/migrate history; fi'

## -----------------------------
## Help
## -----------------------------

help:
	@echo "Usage: make [target] [ENV=dev|prod]"
	@echo "All targets run through Docker. Local Go/Node is not required."
	@echo ""
	@echo "Targets:"
	@echo "  init            Initialize project identifiers (defaults to directory name)"
	@echo "  up              Start containers (default: dev)"
	@echo "  build           Build containers"
	@echo "  build_no_cache  Build containers without cache"
	@echo "  build_prod      Build the production API image"
	@echo "  down            Stop and remove containers and networks"
	@echo "  down_volumes    Stop and remove containers, networks, and volumes"
	@echo "  stop            Stop containers only"
	@echo "  exec            Enter api container shell"
	@echo "  shell           Start a one-off api shell"
	@echo "  logs            Show api logs"
	@echo "  ps              Show container status"
	@echo "  reup            Restart environment (down + up)"
	@echo "  check           Run format checks, Go vet, and tests"
	@echo "  lint            Run Go vet inside the api container"
	@echo "  format          Format Go files with gofmt"
	@echo "  format_check    Check Go formatting"
	@echo "  test            Run Go tests inside the api container"
	@echo "  test_e2e        Run the full HTTP API contract in isolation"
	@echo "  smoke           Call /health from the running api container"
	@echo "  routes          Print Gin route paths from the api container"
	@echo "  migrate         Run database migrations"
	@echo "  current         Show current migration"
	@echo "  history         Show applied migration history"
	@echo ""
	@echo "Examples:"
	@echo "  make build"
	@echo "  make up"
	@echo "  make migrate"
	@echo "  make smoke"
	@echo "  make build ENV=prod"
