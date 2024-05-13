include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

local-migration-create:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} create init sql

generate:
	make generate-chat-server-api
	make generate-access-api

generate-chat-server-api:
	mkdir -p pkg/chat-server_v1
	protoc --proto_path=api/chat-server_v1 \
		--go_out=pkg/chat-server_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go \
		--go-grpc_out=pkg/chat-server_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
		api/chat-server_v1/chat-server.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 \
		--go_out=pkg/access_v1 --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=./bin/protoc-gen-go \
		--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
		api/access_v1/access.proto

build:
	GOOS=linux GOARCH=amd64 go build -o chat_server cmd/server/main.go

test:
	go clean -testcache
	go test github.com/semho/chat-microservices/chat-server/internal/service/... \
			github.com/semho/chat-microservices/chat-server/internal/api/... -covermode count -count 5


test-coverage:
	go clean -testcache
	go test github.com/semho/chat-microservices/chat-server/internal/service/... \
            github.com/semho/chat-microservices/chat-server/internal/api/... -covermode count -coverprofile=coverage.tmp.out -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore