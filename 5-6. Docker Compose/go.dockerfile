FROM golang:1.26.1-alpine

WORKDIR /app

COPY ./.env .
COPY ./go.mod .
COPY ./go.sum .

COPY ./cmd ./cmd

RUN CGO_ENABLED=0 go build -o main ./cmd/web

CMD ["./main"]