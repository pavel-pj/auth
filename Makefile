include .env

LOCAL_BIN:=$(CURDIR)/bin
LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(POSTGRES_DB) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD)"


install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.33.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.33.0
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7



generate: 
	make generate-user-api 
	make generate-auth-api 
	make generate-access-api
	make generate-openapi
	mkdir -p pkg/swagger
	$(LOCAL_BIN)/statik -src=pkg/swagger -dest=pkg -include="*.css,*.html,*.js,*.json,*.png" -f

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 --proto_path vendor.protogen \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/user_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	api/user_v1/user.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc --proto_path api/auth_v1 --proto_path vendor.protogen \
	--go_out=pkg/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/auth_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	api/auth_v1/auth.proto

generate-access-api:
	mkdir -p pkg/access_v1
	protoc --proto_path api/access_v1 --proto_path vendor.protogen \
	--go_out=pkg/access_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/access_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--grpc-gateway_out=pkg/access_v1 --grpc-gateway_opt=paths=source_relative \
	--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
	api/access_v1/access.proto	

# OpenAPI is generated once, in a single protoc invocation, so allow_merge
# actually merges user_v1 and auth_v1 into one api.swagger.json.
generate-openapi: generate-user-api generate-auth-api
	mkdir -p pkg/swagger
	protoc --proto_path api --proto_path vendor.protogen \
	--openapiv2_out=allow_merge=true,merge_file_name=api:pkg/swagger \
	--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
	api/user_v1/user.proto api/auth_v1/auth.proto


server-up:
	go run cmd/grpc_server/main.go

client-up:
	go run cmd/grpc_client/main.go	


goose-status:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) status -v
goose-up:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) up -v

goose-down:
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) postgres $(LOCAL_MIGRATION_DSN) down -v

goose-create:
	mkdir -p $(LOCAL_MIGRATION_DIR)
	$(LOCAL_BIN)/goose -dir $(LOCAL_MIGRATION_DIR) create $(NAME) sql

lint:
	$(LOCAL_BIN)/golangci-lint run ./...


 
tests-cover:	
	go test -count=5 \
		-coverpkg=auth/internal/api/user,auth/internal/service/user \
		-coverprofile=coverage.out \
		./internal/api/user/tests ./internal/service/user/tests
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html	
 

lib:
	GOPROXY=direct go get github.com/pavel-pj/platform_lib@latest


vendor-proto:
	@if [ ! -f vendor.protogen/buf/validate/validate.proto ]; then \
		rm -rf vendor.protogen/protovalidate &&\
		git clone --depth=1 https://github.com/bufbuild/protovalidate vendor.protogen/protovalidate &&\
		mkdir -p vendor.protogen/buf/validate &&\
		cp vendor.protogen/protovalidate/proto/protovalidate/buf/validate/*.proto vendor.protogen/buf/validate/ &&\
		rm -rf vendor.protogen/protovalidate ;\
	fi	
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
		mkdir -p vendor.protogen/google/ &&\
		mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
		rm -rf vendor.protogen/googleapis ;\
	fi
	@if [ ! -f vendor.protogen/protoc-gen-openapiv2/options/annotations.proto ]; then \
		rm -rf vendor.protogen/openapiv2 &&\
		git clone --depth=1 --branch v2.33.0 https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
		mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
		mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
		rm -rf vendor.protogen/openapiv2 ;\
	fi

jwt-string:
	openssl rand -base64 48


gen-cert:
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext
