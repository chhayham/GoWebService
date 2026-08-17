# Build stage
FROM golang:1.26.6 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o webservice main.go

# Final runtime stage
FROM alpine:latest

# Install ca-certificates in case the API needs to make HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/webservice .

# Copy the static folder (CSS/assets)
COPY --from=builder /app/static ./static

# Expose port 8080 as defined in main.go
EXPOSE 8080

# Run the application
CMD ["./webservice"]