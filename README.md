# ms-catalog

[![CI](https://github.com/selfshop-dev/ms-catalog/actions/workflows/ci.yml/badge.svg)](https://github.com/selfshop-dev/ms-catalog/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/selfshop-dev/ms-catalog/branch/main/graph/badge.svg)](https://codecov.io/gh/selfshop-dev/ms-catalog)
[![Go Report Card](https://goreportcard.com/badge/github.com/selfshop-dev/ms-catalog)](https://goreportcard.com/report/github.com/selfshop-dev/ms-catalog)
[![Go version](https://img.shields.io/github/go-mod/go-version/selfshop-dev/ms-catalog)](go.mod)
[![License](https://img.shields.io/github/license/selfshop-dev/ms-catalog)](LICENSE)

REST-микросервис управления каталогом товаров на Go. Часть платформы [selfshop-dev](https://github.com/selfshop-dev).

## Обзор

`ms-catalog` реализует полный жизненный цикл товара: создание, обновление, управление статусом, мягкое удаление и постраничный просмотр активных товаров. Сервис построен по принципу чистой архитектуры с явным разделением на слои domain, usecase, handler и storage.

Подробнее об архитектуре — в [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).  
Полная спецификация API — в [docs/API.md](docs/API.md).  
Обоснование технологических решений — в [docs/DECISIONS](docs/DECISIONS).

### Возможности

- CRUD для товаров с поддержкой статусов `active`, `inactive`, `draft`, `archived`
- Поиск по UUID и slug
- Пагинированный список активных товаров
- Валидация на уровне домена и на уровне HTTP-запроса
- RFC 9457 problem details для всех ошибочных ответов
- Структурированное логирование через zap с runtime-управлением уровнем
- Health-эндпоинты `/health/alive` и `/health/ready`
- Swagger UI в dev-режиме на `/swagger`
- Hot-reload в dev-окружении через air

### Быстрый старт

```bash
# Клонировать репозиторий
git clone https://github.com/selfshop-dev/ms-catalog
cd ms-catalog

# Скопировать и заполнить переменные окружения
cp .env.example .env

# Запустить dev-окружение с hot-reload
make dev

# Применить миграции
make mig-apply

# Открыть Swagger UI
open http://localhost:8080/swagger
```

## Конфигурация

Сервис конфигурируется через переменные окружения с префиксом `INIT__`. Двойное подчёркивание `__` — разделитель вложенности (`INIT__DB__DSN` → `db.dsn`).

| Переменная | По умолчанию | Описание |
|---|---|---|
| `INIT__APP__NAME` | `my.ms-todo` | Имя сервиса в логах |
| `INIT__APP__RUNMODE` | `prod` | Режим: `dev` или `prod` |
| `INIT__LOG__MIN_LEVEL` | `info` | Уровень логирования |
| `INIT__LOG__FORMAT` | `auto` | Формат: `json`, `console`, `auto` |
| `INIT__ENTRY__HTTP__PORT` | — | TCP-порт сервера (обязательно) |
| `INIT__ENTRY__HTTP__READ_TIMEOUT` | `10s` | Таймаут чтения запроса |
| `INIT__ENTRY__HTTP__REQUEST_TIMEOUT` | `15s` | Дедлайн контекста обработчика |
| `INIT__ENTRY__HTTP__WRITE_TIMEOUT` | `20s` | Таймаут записи ответа |
| `INIT__ENTRY__HTTP__IDLE_TIMEOUT` | `90s` | Keep-alive таймаут |
| `INIT__DB__DSN` | — | PostgreSQL DSN (обязательно) |
| `INIT__DB__MAX_CONNS` | `10` | Максимальный размер пула |
| `INIT__DB__MIN_CONNS` | `5` | Минимальный размер пула |
| `INIT__DB__MAX_CONN_LIFETIME` | `5m` | Максимальное время жизни соединения |
| `INIT__DB__MAX_CONN_IDLE_TIME` | `30s` | Максимальное время простоя соединения |

### Разработка

Подробнее о рабочем процессе и правилах оформления коммитов — в [CONTRIBUTING.md](CONTRIBUTING.md).

```bash
make test     # unit-тесты с покрытием
make test-e2e # end-to-end тесты (требует Docker)

make lint # статический анализ (golangci-lint)

make sqlc-gen # регенерировать Go-код из SQL
make swag-gen # регенерировать OpenAPI spec
make code-gen # go generate ./...
```

## Лицензия

[`MIT`](LICENSE) © 2026-present [`selfshop-dev`](https://github.com/selfshop-dev)