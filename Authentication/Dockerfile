# FROM golang:alpine AS builder
# WORKDIR /app

# COPY ./AuthenticationService/go.mod .
# COPY ./AuthenticationService/go.sum .
# RUN go mod download

# COPY ./AuthenticationService/ .
# RUN go build -o auth-service

# FROM alpine
# COPY --from=builder /app/auth-service /usr/bin
# EXPOSE 3001
# ENTRYPOINT ["go", "run", "main.go"]


# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./Authentication/go.mod ./Authentication/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./Authentication/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 3001 to the outside world
EXPOSE 3001

# Command to run the executable
CMD ["./main"]