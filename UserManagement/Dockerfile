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

FROM golang:latest as builder

WORKDIR /app

COPY ./UserManagement/go.mod ./UserManagement/go.sum ./

RUN go mod download

COPY ./UserManagement/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3000

CMD ["./main"]