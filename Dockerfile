# Use the official Golang image as the base image
FROM golang:1.19-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the container, and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code to the container
COPY . .

# Build the Go application
RUN go build -o taskmanager main.go

# Expose port 8080 so that the app is accessible
EXPOSE 8080

# Run the task manager binary
CMD ["./taskmanager"]
