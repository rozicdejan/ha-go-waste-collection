# Use a minimal Go image to run the Go application
FROM golang:1.18-alpine

# Set the working directory
WORKDIR /app

# Copy the Go script into the container
COPY main.go .

# Install dependencies and build the Go script
RUN go mod init waste-collection && \
    go mod tidy && \
    go build -o waste-collection main.go

# Add a run script (see below)
COPY run.sh /run.sh
RUN chmod +x /run.sh

# Run the script when the container starts
CMD ["/run.sh"]
