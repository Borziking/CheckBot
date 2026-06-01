FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bot .
COPY Roboto-Regular.ttf .
COPY Roboto-Medium.ttf .
COPY Roboto-Bold.ttf .
COPY config.json .

CMD ["./bot"]