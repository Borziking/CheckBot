FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bot .
COPY assets ./assets
# config.json в /app — это дефолт-сид; рабочая копия живёт в DATA_DIR.
COPY config.json .

# Изменяемые данные (config.json, users.json) хранятся в томе и переживают
# пересоздание контейнера.
ENV DATA_DIR=/data
VOLUME ["/data"]

CMD ["./bot"]