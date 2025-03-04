FROM golang:1.18-alpine as builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o idp ./cmd/idp

# Use a small alpine image
FROM alpine:3.15

WORKDIR /app

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/idp .

# Copy config and web directories
COPY --from=builder /app/config ./config
COPY --from=builder /app/web ./web

# Expose the service port
EXPOSE 8080

# Run the service
CMD ["./idp"]