# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY ./src ./src

COPY go.mod go.mod
COPY go.sum go.sum

RUN go test -v ./src/...

RUN go build -o main ./src

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE ${PORT}

CMD ["./main"]