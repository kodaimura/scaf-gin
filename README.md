# scaf-gin-api

Gin backend scaffold.

This template is intended to run through Docker. Local Go and Node are not
required for normal development.
PostgreSQL is the supported database. Configure it with `DATABASE_URL`.
The server stores and emits timestamps in UTC.

## Development

```sh
cp .env.example .env
make build
make up
make migrate
```

Useful commands:

```sh
make logs
make exec
make shell
make check
make test
make smoke
make routes
make migrate
make current
make history
make down_volumes
```

The API runs at `http://localhost:8000/api`.
Health check is available at `http://localhost:8000/health`.
MailHog is available at `http://localhost:8025`.
Host ports are bound to `127.0.0.1` by default. Set `API_BIND_HOST=0.0.0.0`
only when the API must be reachable from outside the host.
Source changes are reloaded automatically in the development container.

## Structure

```text
cmd/          # application entrypoints
config/       # environment loading and validation
internal/
  app/        # Gin bootstrap, middleware, and routes
  core/       # config-backed infrastructure: auth, database, logger, mailer
  handler/    # HTTP request/response handling
  model/      # database models
  module/     # persistence-oriented domain modules
  usecase/    # application use cases
migrations/   # numbered PostgreSQL SQL migrations
```

Migration files are numbered sequentially from `001`.
Applied migrations are recorded in the `schema_migrations` table.

Use production compose settings with `ENV=prod`.

```sh
cp .env.example .env
# Edit production secrets and DATABASE_URL in .env.
make build ENV=prod
make migrate ENV=prod
make up ENV=prod
```

The development database is stored in the Docker named volume
`scaf-gin-api_db_data`.
