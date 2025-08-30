# Pake debian at least pake 2GB storage, im not that rich
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files (if they exist)
COPY go.mod ./

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main
COPY .env .env
RUN go run migrate/migrate.go
# Use a small alpine image for the final container
FROM alpine:3.21

# Add CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

# Expose port 8080
EXPOSE 3000

CMD ["./main"]