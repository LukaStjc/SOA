# FROM golang:alpine AS builder
# WORKDIR /app

# COPY ./UserService/go.mod .
# COPY ./UserService/go.sum .
# RUN go mod download

# COPY ./UserService/ .
# RUN go build -o user-service

# FROM alpine
# COPY --from=builder /app/user-service /usr/bin
# EXPOSE 3000
# ENTRYPOINT ["go", "run", "main.go"]

# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./UserManagement/go.mod ./UserManagement/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./UserManagement/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]