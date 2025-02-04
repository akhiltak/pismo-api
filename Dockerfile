# Start from the official Go image
FROM golang:1.23.5-alpine AS builder

# Install make
RUN apk add --no-cache make

# Set the working directory inside the container
WORKDIR /app

# Copy only go.mod and go.sum (to cache dependencies)
COPY go.mod go.sum Makefile ./

# Install dependencies (this will only run when go.mod/go.sum change)
RUN make deps

# Now copy the rest of the project
COPY . .

# Set BUILD_DIR environment variable
ENV BUILD_DIR=build/

# Run make commands to build the application
RUN make build

# Start a new stage from scratch to reduce image size
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/build/pismo-backend .
COPY db/migrations/ db/migrations/


# Expose the port the app runs on
EXPOSE 2090

# Command to run the executable
CMD ["./pismo-backend"]
