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
# Seed config copied into the image; the working copy lives in DATA_DIR.
COPY config.example.json ./config.json

# Mutable data (config.json, users.json) lives in a volume and survives recreation.
ENV DATA_DIR=/data
# Where render looks for fonts (assets/fonts), independent of the working directory.
ENV ASSETS_DIR=/app/assets
VOLUME ["/data"]

CMD ["./bot"]