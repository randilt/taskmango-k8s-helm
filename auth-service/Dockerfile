# Build stage
FROM golang:1.23.6 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-service .
COPY app.env .

EXPOSE 9090 9091

CMD ["./auth-service"]