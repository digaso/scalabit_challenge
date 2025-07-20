# Build Stage
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila o main que está em cmd/scalabit
RUN go build -o scalabit-api ./cmd

# Runtime Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/scalabit-api .

EXPOSE 8080

CMD ["./scalabit-api"]
