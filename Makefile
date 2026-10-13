include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD)"


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

generate: generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

server-up:
	go run internal/grpc_server/main.go


goose-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v
goose-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

goose-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

goose-create:
	mkdir -p $(LOCAL_MIGRATION_DIR)
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create $(NAME) sql
 