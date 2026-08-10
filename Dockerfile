FROM golang:1.24 AS dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@v1.62.0

CMD ["air", "-c", ".air.toml"]


FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/api ./cmd
RUN go build -o /out/healthcheck ./cmd/healthcheck


FROM debian:bookworm-slim AS runtime

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/api /app/api
COPY --from=builder /out/healthcheck /app/healthcheck

EXPOSE 8000

USER 65532:65532

CMD ["/app/api"]
