# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Install necessary build dependencies for CGO (sqlite3 needs them)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

# Set CGO_ENABLED=1 to enable CGO
ENV CGO_ENABLED=1

RUN go build -o main .

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder image
COPY --from=builder /app/main .

# Install the SQLite library in the runtime container
RUN apk add --no-cache sqlite

RUN mkdir -p /app/data

EXPOSE ${PORT}

CMD ["./main"]
