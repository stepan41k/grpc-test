FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main cmd/server/main.go


FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

CMD ["./main"]