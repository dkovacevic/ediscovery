# Use the same Debian version for both build and runtime stages
FROM golang:1.21 AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
RUN mkdir src

COPY src/* ./src/

# Build the Go app
RUN go build -o lh-whatsapp ./src

# Second stage: create the runtime image
FROM debian:bullseye-slim

RUN mkdir /opt/whatsapp

# Set the current working directory inside the container
WORKDIR /opt/whatsapp

# Install necessary packages including the `file` command for debugging
RUN apt-get update && apt-get install -y --no-install-recommends sqlite3 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/lh-whatsapp .

RUN chmod +x lh-whatsapp

# Command to run the executable
CMD ["./lh-whatsapp"]
