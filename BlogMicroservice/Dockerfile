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


FROM golang:latest as builder

WORKDIR /app

COPY ./BlogMicroservice/go.mod ./BlogMicroservice/go.sum ./

RUN go mod download

COPY ./BlogMicroservice/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]