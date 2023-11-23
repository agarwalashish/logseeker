# Build Stage
FROM golang:1.21 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
COPY logs/ /app/logs

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o logseeker .

# Final Stage
FROM alpine:latest

WORKDIR /root/

# Import the Certificate-Authority certificates for enabling HTTPS
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/logseeker .

# Copy the sample log files into /var/log
COPY --from=builder /app/logs /var/log


# Command to run the executable
ENTRYPOINT ["./logseeker"]
