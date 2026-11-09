# auth — сервис пользователей (gRPC + REST + Swagger)

Микросервис авторизации/пользователей на Go: gRPC-сервер, HTTP/JSON через `grpc-gateway`, валидация через Protovalidate, PostgreSQL, миграции через `goose`, Swagger UI, встроенный в бинарник через `statik`.

## Возможности

- CRUD пользователей (`Store`, `Get`, `Index`, `Update`, `Delete`).
- gRPC и REST/JSON из одного описания (`api/user_v1/user.proto`).
- Валидация входных сообщений на уровне gRPC-интерсептора (Protovalidate, CEL-совместимая схема).
- Swagger UI и OpenAPI 2.0 спецификация.
- PostgreSQL + миграции `goose`.
- Транзакционная запись: пользователь + запись в `logs`.
- Docker-образы приложения и мигратора, CI/CD через GitHub Actions.
- Юнит-тесты (minimock, testify), покрытие, golangci-lint.

## Стек

- Go 1.27+
- gRPC, `grpc-gateway/v2`
- Protovalidate (`buf.build/go/protovalidate`) + `buf/validate`
- PostgreSQL 18, `pgx/v5`, `squirrel`, `platform_lib`
- `goose` (миграции), `statik` (встраивание Swagger UI)
- `minimock`, `testify`, `golangci-lint`
- Docker / Docker Compose, GitHub Actions

## Архитектура и порты

Один бинарник поднимает три сервера:

| Компонент | Порт по умолчанию | Переменная |
|---|---|---|
| gRPC | `50560` | `SERVER_PORT` |
| HTTP/JSON gateway | `50570` | `SERVER_HTTP_PORT` |
| Swagger UI | `50580` | `SERVER_SWAGGER_PORT` |
| PostgreSQL (host) | `5091` → `5432` | `PG_PORT` |

Путь запроса:

```text
Swagger UI / HTTP client
        │
        ▼
HTTP gateway (50570) ──► gRPC (50560) ──► service ──► repository ──► PostgreSQL
```

Валидация Protovalidate выполняется в gRPC-интерсепторе **до** вызова хендлера.

## Структура проекта

```text
.
├── api/user_v1/user.proto        # контракт: сервис, сообщения, HTTP-аннотации, правила валидации
├── cmd/grpc_server/main.go       # точка входа
├── internal/
│   ├── api/user/                 # gRPC-реализация (Store/Get/Index/Update/Delete)
│   ├── app/                      # сборка приложения, запуск серверов, DI через serviceProvider
│   ├── config/                   # конфиги и чтение env
│   ├── converter/                # proto <-> model
│   ├── interceptor/              # recovery + validation интерсепторы
│   ├── model/                    # доменные модели
│   ├── repository/               # доступ к БД (user, log)
│   └── service/                  # бизнес-логика, транзакции
├── migrations/                   # goose-миграции
├── pkg/user_v1/                  # сгенерированный Go-код из proto
├── pkg/swagger/                  # Swagger UI + api.swagger.json
├── pkg/statik/                   # сгенерированный встроенный FS
├── vendor.protogen/              # вендор .proto (buf/validate, google/api, openapiv2) — не в git
├── compose.dev.yaml              # PostgreSQL для локальной разработки
├── compose.prod.yaml             # production: postgres, migrator, grpc-server
├── Dockerfile                    # образ приложения
├── migration.Dockerfile          # образ мигратора
└── Makefile
```

## Требования

- Go 1.27+
- `protoc` (в PATH)
- `make`, `git`, `docker`, `docker compose`
- Доступ к GitHub для `make vendor-proto`
- `npm` не нужен: Swagger UI уже лежит в `pkg/swagger`

## Быстрый старт

> `Makefile` делает `include .env`, поэтому `.env` должен существовать **до** запуска любых `make`-целей.

1. Создать `.env` (пример ниже).

2. Поднять PostgreSQL:

```bash
docker compose -f compose.dev.yaml up -d postgres
```

3. Установить инструменты в `./bin` (protoc-плагины, goose, golangci-lint, statik):

```bash
make install-deps
```

Это устанавливает только CLI-утилиты. Go-библиотеки приложения тянутся модулями:

```bash
go mod download
```

4. Подтянуть `.proto`-зависимости (buf/validate, google/api, openapiv2 options):

```bash
make vendor-proto
```

Повторный запуск безопасен: если файлы уже есть, цель молча ничего не делает.

5. Сгенерировать код и Swagger-ассеты:

```bash
make generate
```

6. Применить миграции:

```bash
make goose-up
```

7. Запустить сервис:

```bash
make server-up
```

После старта:

- gRPC — `localhost:50560`
- HTTP API — `http://localhost:50570`
- Swagger UI — `http://localhost:50580/`

## Конфигурация

Пример `.env` (значения-заглушки, реальные секреты не коммитить):

```dotenv
# --- PostgreSQL ---
POSTGRES_DB=db_auth
POSTGRES_USER=auth
POSTGRES_PASSWORD=change_me
DB_HOST=localhost
PG_PORT=5091        # host-порт из compose.dev
DB_PORT=5091        # порт, к которому подключается приложение
MIGRATION_DIR=./migrations

# --- gRPC ---
SERVER_HOST=
SERVER_PORT=50560

# --- HTTP gateway ---
SERVER_HTTP_HOST=
SERVER_HTTP_PORT=50570

# --- Swagger ---
SERVER_SWAGGER_HOST=localhost
SERVER_SWAGGER_PORT=50580
```

Пояснения:

- `POSTGRES_*` используются и приложением, и контейнером Postgres.
- `PG_PORT` — порт на хосте (compose и `goose`), `DB_PORT` — порт, который использует само приложение.
- `SERVER_*_HOST` пустой означает «слушать на всех интерфейсах» (кроме Swagger, там по умолчанию `localhost`).

Для production дополнительно задаётся:

```dotenv
DOCKERHUB_USERNAME=<ваш-dockerhub-username>
```

## API

Базовый URL: `http://localhost:50570`.

| Метод | Путь | Описание |
|---|---|---|
| `POST` | `/user/v1/create` | Создать пользователя |
| `GET` | `/user/v1/list?offset=0&limit=10` | Список пользователей |
| `GET` | `/user/v1/{id}` | Пользователь по ID |
| `PATCH` | `/user/v1/{id}` | Обновить пользователя |
| `DELETE` | `/user/v1/{id}` | Удалить пользователя |

JSON-имена — camelCase (`passwordConfirm`, `createdAt`, `updatedAt`).

### Создание

```bash
curl -X POST http://localhost:50570/user/v1/create \
  -H 'Content-Type: application/json' \
  -d '{
    "user": {
      "info": {
        "name": "Ivan",
        "email": "ivan@example.com",
        "roles": "USER"
      },
      "password": "secret123",
      "passwordConfirm": "secret123"
    }
  }'
```

### Список

```bash
curl 'http://localhost:50570/user/v1/list?offset=0&limit=10'
```

### Получение

```bash
curl http://localhost:50570/user/v1/1
```

### Обновление

```bash
curl -X PATCH http://localhost:50570/user/v1/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Ivan2",
    "email": "ivan2@example.com",
    "roles": "ADMIN"
  }'
```

### Удаление

```bash
curl -X DELETE http://localhost:50570/user/v1/1
```

### Валидация

Правила задаются в `api/user_v1/user.proto` через `buf.validate`:

- `name` — от 1 до 7 символов;
- `email` — валидный email, не длиннее 70 символов;
- в `UpdateRequest.name` — те же ограничения.

Валидация выполняется в `internal/interceptor/validate.go` (Protovalidate) и возвращает `codes.InvalidArgument` → HTTP `400`.

> Важно: `passwordConfirm` — именно camelCase. Если отправить `password_confirm`, поле будет проигнорировано, пароли «не совпадут», и вернётся ошибка.

### Ошибки

- `400` — не прошла валидация Protovalidate.
- `500` — ошибки БД (например, `users.email UNIQUE` → `duplicate key`), несовпадение паролей и прочие внутренние ошибки.

## Swagger UI

- UI: `http://localhost:50580/`
- Спецификация: `http://localhost:50580/api.swagger.json`
- `host` в спецификации: `localhost:50570` (см. `openapiv2_swagger` в `user.proto`)

Swagger UI — это статика из `pkg/swagger`, встроенная в бинарник через `statik`:

- `make generate` генерирует `pkg/swagger/api.swagger.json` и `pkg/statik/statik.go`;
- в `internal/app/app.go` обязателен blank-import:

```go
_ "auth/pkg/statik"
```

Без него при старте будет ошибка:

```text
statik/fs: no zip data registered
```

Обновление Swagger UI:

1. положить новые файлы в `pkg/swagger` (Swagger UI dist + свой `index.html`/`swagger-initializer.js`);
2. `make generate` (перегенерирует statik);
3. перезапустить сервис.

CORS для UI настроен на HTTP-сервере (`50570`), потому что страница на `50580` обращается к API на другом origin. В dev сейчас разрешены все origin'ы; для production ограничьте список.

## Генерация кода

```bash
make vendor-proto   # .proto-зависимости (один раз / при изменении версий)
make generate       # protoc + openapiv2 + statik
```

`make generate` выполняет:

- `protoc-gen-go`, `protoc-gen-go-grpc` → `pkg/user_v1`;
- `protoc-gen-grpc-gateway` → HTTP-хендлеры;
- `protoc-gen-openapiv2` → `pkg/swagger/api.swagger.json`;
- `statik` → `pkg/statik/statik.go`.

После изменения `user.proto` или ассетов Swagger повторно запускайте `make generate`.

## Миграции

```bash
make goose-status    # статус
make goose-up        # применить
make goose-down      # откатить
make goose-create NAME=add_something
```

DSN собирается из `.env` (`PG_PORT`, `POSTGRES_*`, `MIGRATION_DIR`).

В production миграции выполняет отдельный контейнер `migrator` (`migration.Dockerfile` + `migration.sh`), который ждёт Postgres и запускает `goose up`.

## Тесты и линт

```bash
go test ./...
make tests-cover     # 5 прогонов + coverage.out/coverage.html
make lint            # golangci-lint
```

## Docker и деплой

Локальная сборка:

```bash
docker build -t <user>/auth:latest .
docker build -f migration.Dockerfile -t <user>/auth-migrator:latest .
```

Production (`compose.prod.yaml`) использует готовые образы:

```bash
DOCKERHUB_USERNAME=<user> docker compose -f compose.prod.yaml up -d
```

Состав: `postgres`, `migrator`, `grpc-server`.

> В `compose.prod.yaml` наружу проброшен только gRPC-порт. HTTP gateway (`50570`) и Swagger (`50580`) доступны через reverse proxy/nginx или после добавления `ports` в compose.

CI/CD (`.github/workflows/deployment.yml`, push в `main`):

1. собирает и пушит `<user>/auth:latest` и `<user>/auth-migrator:latest`;
2. копирует `compose.prod.yaml` на сервер;
3. выполняет `docker compose pull && docker compose up -d`.

Необходимые secrets:

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`
- `SERVER_HOST`
- `SERVER_USER`
- `SSH_PRIVATE_KEY`

## Troubleshooting

**`include .env: No such file or directory`**
Создайте `.env` (см. раздел «Конфигурация»).

**`statik/fs: no zip data registered`**
Нет blank-import `_ "auth/pkg/statik"` или устарел `pkg/statik`. Выполните `make generate` и перезапустите.

**`no required module provides package ...`**
Добавьте зависимость (`go get <module>`) или `go mod tidy`; для CLI-утилит — `make install-deps`.

**`make vendor-proto` ничего не выводит**
Это нормально: все проверки `if [ ! -f ... ]` уже удовлетворены.

**Валидация не срабатывает**
Проверьте, что `internal/interceptor/validate.go` использует `protovalidate.Validate`, а интерсептор зарегистрирован в `initGRPCServer`.

**CORS-ошибка в Swagger**
API-сервер (`50570`) должен отдавать `Access-Control-Allow-Origin` для origin страницы Swagger (`50580`). `localhost` и `127.0.0.1` — разные origin.

**`duplicate key value violates unique constraint "users_email_key"`**
Email уже существует (`users.email UNIQUE`). Используйте другой email.

**`passwords and password confirmation are not equal`**
Проверьте JSON-поле `passwordConfirm` (не `password_confirm`).

**Порт в Dockerfile**
`EXPOSE 55555` в `Dockerfile` — информационный; реальный порт задаётся `SERVER_PORT` (по умолчанию `50560`) и маппингом в `compose.prod.yaml`.

## Полезные зависимости

- DI, БД-клиент, closer: [pavel-pj/platform_lib](https://github.com/pavel-pj/platform_lib)
- Валидация: [protovalidate](https://protovalidate.com/)
- Gateway/OpenAPI: [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)