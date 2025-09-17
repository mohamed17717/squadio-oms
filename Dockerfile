# Start with the official Golang image
FROM golang:1.23-alpine

# Set environment variables
ENV GO111MODULE=on
ENV GOPATH=/go

# Create app directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port your app will run on
EXPOSE 8080

# Run the executable
CMD ["./main"]
