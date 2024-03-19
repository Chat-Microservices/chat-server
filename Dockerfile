FROM golang:1.21.8-alpine AS builder

COPY . /github.com/semho/chat-microservices/chat-server/
WORKDIR /github.com/semho/chat-microservices/chat-server/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/server/main.go

FROM alpine:3.19.1

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/* \

WORKDIR /root/
COPY --from=builder /github.com/semho/chat-microservices/chat-server/bin/chat_server .
COPY --from=builder /github.com/semho/chat-microservices/chat-server/entrypoint.sh .
COPY --from=builder /github.com/semho/chat-microservices/chat-server/migrations ./migrations

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
RUN chmod +x entrypoint.sh