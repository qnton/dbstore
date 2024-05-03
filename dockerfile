# Base image
FROM golang:1.22 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code.
COPY . .

# Build the dumper sub-project
WORKDIR /app/dumper
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/dumper

# Build the main sub-project
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/dumper .
COPY --from=builder /app/bin/main .
CMD ["./main"]