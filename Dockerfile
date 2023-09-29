FROM golang:1.21.0-alpine3.18

RUN mkdir /app

WORKDIR /app

COPY . .

CMD ["go", "test", "-v", "."]