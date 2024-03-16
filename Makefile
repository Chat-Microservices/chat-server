include .env.local

LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

generate:
	make generate-chat-server-api

generate-chat-server-api:
	mkdir -p pkg/chat-server_v1
	protoc --proto_path=api/chat-server_v1 \
		--go_out=pkg/chat-server_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go \
		--go-grpc_out=pkg/chat-server_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
		api/chat-server_v1/chat-server.proto