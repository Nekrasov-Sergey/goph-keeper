.DEFAULT_GOAL := server

LDFLAGS = -ldflags "\
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildVersion=1.0.0 \
	-X 'github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildDate=$$(date "+%Y-%m-%d %H:%M:%S")' \
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildCommit=$$(git rev-parse --short HEAD)"

.PHONY: first-build
first-build: gen build-server build-client-macos

.PHONY: server
server: build-server
	@./build/server-darwin-arm64

.PHONY: client
client: build-client-macos
	@./build/client-darwin-arm64

.PHONY: build-server
build-server:
	@go build ${LDFLAGS} -o ./build/server-darwin-arm64 ./cmd/server/server.go

.PHONY: build-all-client
build-all-client: build-client-linux build-client-macos build-client-windows

.PHONY: build-client-macos
build-client-macos:
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ./build/client-darwin-arm64 ./cmd/client/client.go

.PHONY: build-client-linux
build-client-linux:
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ./build/client-linux-amd64 ./cmd/client/client.go

.PHONY: build-client-windows
build-client-windows:
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ./build/client-windows-amd64.exe ./cmd/client/client.go

.PHONY: gen
gen: proto-gen generate-certs

.PHONY: proto-gen
proto-gen:
	@echo "🧰 Генерация protobuf и gRPC кода..."
	@mkdir -p proto
	@protoc \
    --go_out=. --go_opt=module=github.com/Nekrasov-Sergey/goph-keeper \
    --go-grpc_out=. --go-grpc_opt=module=github.com/Nekrasov-Sergey/goph-keeper \
    api/proto/keeper.proto

generate-certs:
	@mkdir -p certs
	@openssl req -x509 -nodes -days 365 \
	-newkey rsa:4096 \
	-keyout certs/server.key \
	-out certs/server.crt \
	-config certs/cert.conf \
	-extensions req_ext

.PHONY: fmt
fmt:
	@echo "Форматирование проекта..."
	@goimports -w .
