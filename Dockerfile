# Copyright (c) 2022-2024 Ronan LE MEILLAT
# This program is licensed under the AGPLv3 license.
# Use the official Go image as a parent image
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the source code
COPY *.go ./
COPY not-found.html ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o nginx-ad-auth

# Use a minimal alpine image for the final stage
FROM alpine:3.20

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/nginx-ad-auth .

# Expose the default port
EXPOSE 8080

# Run the application
CMD ["./nginx-ad-auth"]
