DOCKER_COMPOSE := docker compose
ENV ?= dev
DOCKER_COMPOSE_FILE := $(if $(filter prod,$(ENV)),-f docker-compose.prod.yml,-f docker-compose.yml)
DOCKER_COMPOSE_CMD := $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE)
API_SERVICE := api
MIGRATE_SERVICE := migrate

.DEFAULT_GOAL := help

.PHONY: up build build_no_cache down down_volumes stop exec shell logs ps reup check smoke routes migrate current history help

## -----------------------------
## Base Commands
## -----------------------------

up:
	$(DOCKER_COMPOSE_CMD) up -d

build:
	$(DOCKER_COMPOSE_CMD) build

build_no_cache:
	$(DOCKER_COMPOSE_CMD) build --no-cache

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

check:
	$(DOCKER_COMPOSE_CMD) run --rm --no-deps $(API_SERVICE) go test ./...

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
	@echo "  up              Start containers (default: dev)"
	@echo "  build           Build containers"
	@echo "  build_no_cache  Build containers without cache"
	@echo "  down            Stop and remove containers and networks"
	@echo "  down_volumes    Stop and remove containers, networks, and volumes"
	@echo "  stop            Stop containers only"
	@echo "  exec            Enter api container shell"
	@echo "  shell           Start a one-off api shell"
	@echo "  logs            Show api logs"
	@echo "  ps              Show container status"
	@echo "  reup            Restart environment (down + up)"
	@echo "  check           Run Go tests inside the api container"
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
