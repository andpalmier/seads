# Stage 1: Build the Go application
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o seads ./cmd/seads/

# Stage 2: Create a minimal runtime container
FROM ghcr.io/go-rod/rod:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage into the minimal runtime container
COPY --from=builder /app/seads /app/seads

# Set the entrypoint for the application
ENTRYPOINT ["/app/seads"]
