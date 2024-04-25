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


FROM golang:latest as builder

WORKDIR /app

COPY ./Authentication/go.mod ./Authentication/go.sum ./

RUN go mod download

COPY ./Authentication/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 3001

CMD ["./main"]