FROM golang:latest

# Set environment variables
ENV PORT=9000

# Set the working directory inside the container
WORKDIR /go/src/user-service

# Copy only the go.mod and go.sum files first (for caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Expose the port for the service
EXPOSE $PORT

# Set the entry point for the container
ENTRYPOINT ["./main"]