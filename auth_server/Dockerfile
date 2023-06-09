# Use the official Go image as the base image
FROM --platform=linux/amd64 golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the common code
COPY . .

# Build the common code (if applicable)
RUN go build -v -o /usr/local/bin/app ./

# No CMD or ENTRYPOINT needed for the common code
CMD ["app"]
