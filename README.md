# Goph Keeper

Сервис для безопасного хранения приватной информации.

## Типы секретов

- Логин и пароль
- Текстовые данные
- Бинарные файлы
- Банковские карты

## Архитектура

- **gRPC сервер** — хранение и шифрование данных
- **CLI клиент** — консольный интерфейс для работы с секретами

---

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone https://github.com/Nekrasov-Sergey/goph-keeper.git
cd goph-keeper
```

### 2. Проверить зависимости

```bash
make deps
```

### 3. Запуск PostgreSQL через Docker

```bash
make docker-up
```

### 4. Настроить конфигурацию

```bash
cp .env.example .env
cp config/local_example.yml config/local.yml
```

Отредактируйте файлы, заменив `change_me` на реальные значения.

### 5. Первичная сборка

```bash
make first-build
```

---

## Запуск

### Сервер

```bash
make server
```

### Клиент

```bash
make client
```

---

## Команды Makefile

| Команда | Описание |
|---------|----------|
| `make help` | Показать справку |
| `make first-build` | Первичная сборка (генерация + компиляция) |
| `make server` | Собрать и запустить сервер |
| `make client` | Собрать и запустить клиент |
| `make build-all-client` | Собрать клиент для всех платформ |
| `make gen` | Генерация кода (protobuf, сертификаты, моки) |
| `make test` | Запуск тестов |
| `make test-cover` | Тесты с отчётом покрытия |
| `make lint` | Запуск линтера |
| `make fmt` | Форматирование кода |
| `make clean` | Удалить артефакты сборки |
| `make docker-up` | Запуск PostgreSQL в Docker |
| `make docker-down` | Остановка PostgreSQL |

---

## Платформы клиента

Клиент собирается для:

| Платформа | Архитектура | Файл |
|-----------|-------------|------|
| macOS | arm64 | `client-darwin-arm64` |
| Linux | amd64 | `client-linux-amd64` |
| Windows | amd64 | `client-windows-amd64.exe` |

---

## Безопасность

| Механизм | Назначение |
|----------|------------|
| TLS | Защита gRPC соединения |
| AES-256 | Шифрование секретов |
| bcrypt | Хэширование паролей пользователей |
| JWT | Аутентификация и авторизация |

Каждый пользователь имеет уникальный ключ шифрования, который хранится в базе данных в зашифрованном виде.

---

## Технологии

- Go 1.21+
- gRPC и Protocol Buffers
- PostgreSQL
- Docker
- AES-256-GCM
- bcrypt
- JWT
- PromptUI
