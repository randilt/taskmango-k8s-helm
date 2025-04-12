# Build stage
FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api-service .
COPY app.env .

EXPOSE 8080 8081

CMD ["./api-service"]