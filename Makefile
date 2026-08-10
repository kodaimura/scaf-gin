DOCKER_COMPOSE = docker compose
COMPOSE_ENV ?= dev
DOCKER_COMPOSE_FILE = $(if $(filter prod production,$(COMPOSE_ENV)),-f docker-compose.prod.yml,)
DOCKER_COMPOSE_CMD = $(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE)
SMOKE_BASE_URL ?= http://localhost:8000

.PHONY: up build down down_volumes stop in indb log ps reup check smoke routes help

up:
	$(DOCKER_COMPOSE_CMD) up -d

build:
	$(DOCKER_COMPOSE_CMD) build --no-cache

down:
	$(DOCKER_COMPOSE_CMD) down

down_volumes:
	$(DOCKER_COMPOSE_CMD) down -v

stop:
	$(DOCKER_COMPOSE_CMD) stop

in:
	$(DOCKER_COMPOSE_CMD) exec api sh

indb:
	$(DOCKER_COMPOSE_CMD) exec db sh

log:
	$(DOCKER_COMPOSE_CMD) logs -f

ps:
	$(DOCKER_COMPOSE_CMD) ps

reup: down up

check:
	$(DOCKER_COMPOSE_CMD) run --rm --no-deps api go test ./...

smoke:
	sh scripts/smoke.sh "$(SMOKE_BASE_URL)"

routes:
	$(DOCKER_COMPOSE_CMD) run --rm -e PRINT_ROUTES=true api

help:
	@echo "Usage: make [target] [COMPOSE_ENV=dev|production]"
	@echo ""
	@echo "Targets:"
	@echo "  up             Start containers"
	@echo "  build          Build containers without cache"
	@echo "  down           Stop and remove containers and networks"
	@echo "  down_volumes   Stop and remove containers, networks, and volumes"
	@echo "  stop           Stop containers"
	@echo "  in             Open a shell in the api container"
	@echo "  indb           Open a shell in the db container"
	@echo "  log            Follow container logs"
	@echo "  ps             Show container status"
	@echo "  reup           Restart containers"
	@echo "  check          Run Go tests in Docker"
	@echo "  smoke          Check /health on a running API"
	@echo "  routes         Print registered routes"
