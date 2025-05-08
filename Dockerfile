FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -v ./...

RUN go build -o kaspi-api-wrapper ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kaspi-api-wrapper .

COPY .env.example .env

EXPOSE 8080

CMD ["./kaspi-api-wrapper"]
