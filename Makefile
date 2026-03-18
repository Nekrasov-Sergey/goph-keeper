.DEFAULT_GOAL := help

# Версия проекта
VERSION ?= 1.0.0

# Флаги линкера
LDFLAGS = -ldflags "\
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildVersion=$(VERSION) \
	-X 'github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildDate=$$(date "+%Y-%m-%d %H:%M:%S")' \
	-X github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info.buildCommit=$$(git rev-parse --short HEAD)"

# Директория для сборки
BUILD_DIR ?= build

.PHONY: help
help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## Удалить артефакты сборки
	@echo "🗑️  Очистка..."
	@rm -rf $(BUILD_DIR)
	@rm -rf internal/proto
	@rm -rf internal/server/service/mocks

.PHONY: first-build
first-build: gen build-server build-client ## Первичная сборка проекта

.PHONY: build-server
build-server: ## Собрать сервер
	@echo "🔧 Сборка сервера..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/server ./cmd/server

.PHONY: build-client
build-client: ## Собрать клиент под текущую платформу
	@echo "🔧 Сборка клиента..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/client ./cmd/client

.PHONY: build-all-client
build-all-client: build-client-linux build-client-macos build-client-windows ## Собрать клиент для всех платформ

.PHONY: build-client-macos
build-client-macos:
	@echo "🔧 Сборка клиента для macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-darwin-arm64 ./cmd/client

.PHONY: build-client-linux
build-client-linux:
	@echo "🔧 Сборка клиента для Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-linux-amd64 ./cmd/client

.PHONY: build-client-windows
build-client-windows:
	@echo "🔧 Сборка клиента для Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/client-windows-amd64.exe ./cmd/client

.PHONY: server
server: build-server ## Собрать и запустить сервер
	@./$(BUILD_DIR)/server

.PHONY: client
client: build-client ## Собрать и запустить клиент
	@./$(BUILD_DIR)/client

.PHONY: gen
gen: proto-gen generate-certs generate-mocks ## Генерация кода (protobuf, сертификаты, моки)

.PHONY: proto-gen
proto-gen: ## Генерация protobuf и gRPC кода
	@echo "📦 Генерация protobuf и gRPC кода..."
	@rm -rf internal/proto
	@mkdir -p internal/proto
	@protoc \
		--go_out=. --go_opt=module=github.com/Nekrasov-Sergey/goph-keeper \
		--go-grpc_out=. --go-grpc_opt=module=github.com/Nekrasov-Sergey/goph-keeper \
		api/proto/keeper.proto

.PHONY: generate-certs
generate-certs: ## Генерация TLS сертификатов
	@echo "🔐 Генерация TLS сертификатов..."
	@mkdir -p certs
	@openssl req -x509 -nodes -days 365 \
		-newkey rsa:4096 \
		-keyout certs/server.key \
		-out certs/server.crt \
		-config certs/cert.conf \
		-extensions req_ext 2>/dev/null || echo "⚠️  Сертификаты уже существуют или ошибка конфигурации"

.PHONY: generate-mocks
generate-mocks: ## Генерация моков для тестов
	@echo "🎭 Генерация моков..."
	@rm -rf internal/server/service/mocks
	@mkdir -p internal/server/service/mocks
	@go generate ./...

.PHONY: fmt
fmt: ## Форматирование кода
	@echo "✨ Форматирование кода..."
	@go fmt ./...
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "💡 Установите goimports: go install golang.org/x/tools/cmd/goimports@latest"

.PHONY: lint
lint: ## Запуск линтера
	@echo "🔍 Запуск линтера..."
	@golangci-lint run

.PHONY: test
test: ## Запуск тестов
	@echo "🧪 Запуск тестов..."
	@go test -race -cover ./...

.PHONY: test-cover
test-cover: ## Запуск тестов с покрытием
	@echo "🧪 Запуск тестов с покрытием..."
	@go test -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Отчёт сохранён в coverage.html"

.PHONY: docker-up
docker-up: ## Запуск PostgreSQL в Docker
	@echo "🐳 Запуск PostgreSQL..."
	@docker-compose up -d

.PHONY: docker-down
docker-down: ## Остановка PostgreSQL
	@echo "🐳 Остановка PostgreSQL..."
	@docker-compose down

.PHONY: deps
deps: ## Проверка зависимостей
	@echo "📋 Проверка зависимостей..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go не установлен"; exit 1; }
	@command -v protoc >/dev/null 2>&1 || { echo "❌ protoc не установлен"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "⚠️  Docker не установлен (опционально)"; }
	@echo "✅ Все зависимости установлены"
