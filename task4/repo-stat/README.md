# Repo Stat

Микросервисный проект для работы с GitHub-репозиториями и подписками.

Сервисы:
- `api` — HTTP API для клиента.
- `subscriber` — хранение подписок в PostgreSQL и проверка существования репозитория через GitHub API.
- `processor` — промежуточный gRPC слой между API и Collector.
- `collector` — сбор данных по репозиториям из GitHub.
- `postgres` — БД подписок.

## Архитектура запросов

Старый флоу (из прошлых заданий):
- `GET /api/repositories/info?url=...`
- `API -> Processor -> Collector -> GitHub`

Новый флоу по подпискам:
- `POST /subscriptions`
- `DELETE /subscriptions/{owner}/{repo}`
- `GET /subscriptions`
- `GET /subscriptions/info`
- `API -> Subscriber` (CRUD подписок)  
- `API -> Processor -> Collector -> Subscriber -> GitHub` (агрегация info по всем подпискам)

## Технологии

- PostgreSQL + `pgxpool`
- `sqlc` для типобезопасных запросов
- `golang-migrate` для миграций
- gRPC (`protobuf`)
- Gin (HTTP API)

## Быстрый старт

Требования:
- Docker + Docker Compose
- GNU Make

Запуск всего проекта одной командой:

```bash
make up
```

API после запуска доступен на:
- Swagger: `http://localhost:28080/swagger/index.html`

Остановка:

```bash
make down
```