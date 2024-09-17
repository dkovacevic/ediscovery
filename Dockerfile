# Use the same Debian version for both build and runtime stages
FROM golang:1.21-alpine AS builder

RUN apk add build-base

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY src/ src/

# Build the Go app
RUN CGO_ENABLED=1 GOOS=linux go build -o ediscovery ./src

# Second stage: create the runtime image
FROM alpine

RUN mkdir /opt/ediscovery

# Set the current working directory inside the container
WORKDIR /opt/whatsapp

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/ediscovery .

COPY static/ static/

RUN chmod +x ediscovery

# Command to run the executable
CMD ["./ediscovery"]
