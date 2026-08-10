FROM golang:1.24 AS dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD ["go", "run", "./cmd"]


FROM golang:1.24 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/api ./cmd


FROM debian:bookworm-slim AS runtime

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /out/api /app/api

EXPOSE 8000

CMD ["/app/api"]
