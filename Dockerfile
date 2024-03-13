FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/seads ./cmd/seads/

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/seads /app/seads
ENTRYPOINT ["/app/seads"]
