# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy migrations
COPY internal/database/migrations ./internal/database/migrations

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/internal/database/migrations ./internal/database/migrations

EXPOSE 3000

CMD ["/app/main"]
