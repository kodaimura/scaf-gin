# scaf-gin-api

Gin backend scaffold.

This template is intended to run through Docker. Local Go and Node are not
required for normal development.

## Development

```sh
cp .env.example .env
make build
make up
```

Useful commands:

```sh
make logs
make exec
make shell
make check
make smoke
make routes
make down_volumes
```

The API runs at `http://localhost:8000/api`.
Health check is available at `http://localhost:8000/health`.
MailHog is available at `http://localhost:8025`.
Source changes are reloaded automatically in the development container.

Use production compose settings with `ENV=prod`.

```sh
cp .env.example .env
# Edit production secrets and database settings in .env.
make build ENV=prod
make up ENV=prod
```

The development database is stored in the Docker named volume
`scaf-gin-api_db_data`.
