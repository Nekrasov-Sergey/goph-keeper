.DEFAULT_GOAL := server

LDFLAGS = -ldflags "\
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildVersion=1.0.0 \
	-X 'github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildDate=$$(date "+%Y-%m-%d %H:%M:%S")' \
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildCommit=$$(git rev-parse --short HEAD)"

.PHONY: server
server: build-server
	@./cmd/server/server

.PHONY: build-server
build-server:
	@go build ${LDFLAGS} -o ./cmd/server/server ./cmd/server/server.go

.PHONY: client
client: build-client
	@./cmd/client/client

.PHONY: build-client
build-client:
	@go build ${LDFLAGS} -o ./cmd/client/client ./cmd/client/client.go

.PHONY: gen
gen: proto-gen

.PHONY: proto-gen
proto-gen:
	@echo "🧰 Генерация protobuf и gRPC кода..."
	@protoc \
    --go_out=. --go_opt=module=github.com/Nekrasov-Sergey/goph-keeper \
    --go-grpc_out=. --go-grpc_opt=module=github.com/Nekrasov-Sergey/goph-keeper \
    api/proto/keeper.proto

generate-certs:
	mkdir -p certs
	openssl req -x509 -nodes -days 365 \
	-newkey rsa:4096 \
	-keyout certs/server.key \
	-out certs/server.crt \
	-config certs/cert.conf \
	-extensions req_ext
