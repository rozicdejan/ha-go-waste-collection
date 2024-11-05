# Start from a base image with Go
FROM golang:1.18-alpine

# Set the working directory
WORKDIR /app

# Copy the Go source code
COPY waste-collection.go .

# Compile the Go app
RUN go build -o waste-collection waste-collection.go

# Execute the Go application
CMD ["/app/waste-collection"]
