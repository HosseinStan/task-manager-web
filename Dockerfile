# Start with a base image that has Go installed
FROM golang:1.20-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod ./
COPY go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o main .

# Final stage - Create a smaller image for production
FROM alpine:latest
WORKDIR /root/

# Copy the built binary from the build stage
COPY --from=build /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
