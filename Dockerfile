FROM golang:latest

WORKDIR /app
RUN go version

COPY . .

RUN go build -o cmd/app/auth_service ./cmd/app/main.go

EXPOSE 8080

CMD ["./cmd/app/auth_service"]