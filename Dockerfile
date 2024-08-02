# Use the same Debian version for both build and runtime stages
FROM debian:bullseye AS builder

# Install Go and necessary packages for CGO
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    ca-certificates \
    gcc \
    libc6-dev \
    && wget https://golang.org/dl/go1.21.1.linux-arm64.tar.gz \
    && tar -C /usr/local -xzf go1.21.1.linux-arm64.tar.gz \
    && rm go1.21.1.linux-arm64.tar.gz \
    && ln -s /usr/local/go/bin/go /usr/local/bin/go

# Set Go environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH
ENV CGO_ENABLED=1

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY main.go .

# Build the Go app
RUN go build -o main .

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
COPY --from=builder /app/main .
COPY examplestore.db .

RUN chmod +x ./main

# Command to run the executable
CMD ["./main"]