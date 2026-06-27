# Build 
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o logger ./cmd/api

# Run
FROM alpine:3.24.1
WORKDIR /app
COPY --from=builder /app/logger .
EXPOSE 8080
CMD [ "./logger" ]