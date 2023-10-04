FROM golang:latest

WORKDIR /app
RUN go version

COPY . .

ENV ENV_FILE=/app/.env

RUN go build -o cmd/app/auth_service ./cmd/app/main.go

WORKDIR /app/cmd/app

CMD ["./auth_service"]