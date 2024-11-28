# Use Go 1.22 bullseye as base image
FROM golang:1.22-bullseye AS base

# Move to working directory /project
WORKDIR /project

# Copy the go.mod files to the /project directory
COPY go.mod ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
RUN go build

# Document the port that may need to be published
EXPOSE 8081

# Start the application
CMD ["./dealls-dating-app"]
