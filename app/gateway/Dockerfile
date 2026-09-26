FROM golang:1.27 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o gateway ./cmd/gateway

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends i2c-tools \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/gateway .

CMD ["./gateway"]
