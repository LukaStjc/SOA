# FROM golang:alpine AS builder
# WORKDIR /app

# COPY ./BlogService/go.mod .
# COPY ./BlogService/go.sum .
# RUN go mod download

# COPY ./BlogService/ .
# RUN go build -o blog-service

# FROM alpine
# COPY --from=builder /app/blog-service /usr/bin
# EXPOSE 8081
# ENTRYPOINT ["go", "run", "main.go"]


# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ./BlogMicroservice/go.mod ./BlogMicroservice/go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY ./BlogMicroservice/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8081 to the outside world
EXPOSE 8081

# Command to run the executable
CMD ["./main"]