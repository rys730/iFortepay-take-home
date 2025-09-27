FROM golang:1.24.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./db /app/db
COPY ./docs /app/docs
COPY ./infrastructure /app/infrastructure

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

FROM alpine:3.20
RUN apk add --no-cache ca-certificates libc6-compat bash

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /go/bin/goose /usr/local/bin/
COPY db/migrations /migrations

CMD ["./main"]