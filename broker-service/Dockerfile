# Build
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o broker ./cmd/api

# Run
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/broker .
EXPOSE 8080
CMD [ "./broker" ]