.DEFAULT_GOAL := server

LDFLAGS = -ldflags "\
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildVersion=1.0.0 \
	-X 'github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildDate=$$(date "+%Y-%m-%d %H:%M:%S")' \
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildCommit=$$(git rev-parse --short HEAD)"

.PHONY: server
server: build-server
	@./build/server-darwin-arm64

.PHONY: build-server
build-server:
	@go build ${LDFLAGS} -o ./build/server-darwin-arm64 ./cmd/server/server.go

.PHONY: client
client: build-macos
	@./build/client-darwin-arm64

.PHONY: build-all
build-all: build-linux build-macos build-windows

.PHONY: build-macos
build-macos:
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ./build/client-darwin-arm64 ./cmd/client/client.go

.PHONY: build-linux
build-linux:
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ./build/client-linux-amd64 ./cmd/client/client.go

.PHONY: build-windows
build-windows:
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ./build/client-windows-amd64.exe ./cmd/client/client.go

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
