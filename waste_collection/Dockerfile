# Start from a base Go image
FROM golang:1.18-alpine

# Install necessary dependencies
RUN apk add --no-cache bash jq

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the Go source code and build it
COPY main.go .
RUN go build -o waste-collection main.go

# Copy the run.sh file and make it executable
COPY run.sh /app/run.sh
RUN chmod +x /app/run.sh

# Use the run.sh script as the entry point
CMD ["/app/run.sh"]
