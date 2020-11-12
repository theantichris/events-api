# Start from the Go base image for building the application
FROM golang:alpine as builder

# Enable Go Modules
ENV GO111MODULE=on

# Install Git
RUN apk update && apk add --no-cache git

# Set curent working directory
WORKDIR /app

# Copy mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Buid the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main .

# Start a new container to run the application
FROM scratch

# Copy the pre-built binary file
COPY --from=builder /app/bin/main .
COPY .env .

# Run executable
CMD ["./main"]