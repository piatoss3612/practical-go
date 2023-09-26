FROM golang:1.21.0-alpine3.18 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./server ./server

FROM alpine:3.18

RUN mkdir /app

COPY --from=builder /app/server /app

WORKDIR /app

CMD ["./server"]