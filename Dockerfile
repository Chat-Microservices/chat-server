FROM golang:1.21.8-alpine AS builder

COPY . /github.com/semho/chat-microservices/chat-server/
WORKDIR /github.com/semho/chat-microservices/chat-server/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/server/main.go

FROM alpine:3.19.1

WORKDIR /root/

COPY --from=builder /github.com/semho/chat-microservices/chat-server/bin/chat_server .
COPY entrypoint.sh /root/entrypoint.sh

RUN chmod +x entrypoint.sh