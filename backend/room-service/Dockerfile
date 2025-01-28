FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY ./src ./src
COPY go.mod go.mod
COPY go.sum go.sum

ENV CGO_ENABLED=1

RUN go test -v ./src/...

RUN go build -o /app/main ./src

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

RUN apk add --no-cache sqlite

RUN mkdir -p /app/data

EXPOSE ${PORT}

CMD ["./main"]