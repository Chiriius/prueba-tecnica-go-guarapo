# Build stage
FROM golang:1.24.2 AS builder

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Habilita CGO para sqlite3
ENV CGO_ENABLED=1
RUN go build -o /go/bin/service ./api/cmd/main.go

# Final Stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /go/bin/service /app/service
COPY --from=builder /go/src/app/api/docs /app/api/docs

RUN chmod +x /app/service

EXPOSE 8080

CMD ["/app/service"]