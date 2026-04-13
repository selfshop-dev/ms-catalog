# Architecture

## Обзор

`ms-catalog` — stateless REST-сервис. Единственная внешняя зависимость — PostgreSQL. Сервис не имеет собственного кэша и не публикует события.

```
┌─────────────┐     HTTP      ┌──────────────────────────────────────┐     pgx       ┌────────────┐
│   клиент    │ ────────────► │            ms-catalog                │ ────────────► │ PostgreSQL │
└─────────────┘               └──────────────────────────────────────┘               └────────────┘
```

### Структура директорий

```
ms-catalog/
├── cmd/ms-catalog/        # точка входа, DI-граф (uber/fx), значения по умолчанию
├── internal/
│   ├── app/               # chi-роутер, сборка middleware-цепочки
│   ├── container/         # конфигурация приложения (lib-config)
│   ├── domain/            # доменная модель и бизнес-правила
│   ├── handler/           # HTTP-обработчики (transport layer)
│   │   └── v1/            # версионированные эндпоинты
│   ├── usecase/           # use cases (application layer)
│   │   └── v1/
│   ├── db/
│   │   ├── gen/           # сгенерированный sqlc-код (не редактировать вручную)
│   │   ├── dbstorage/     # адаптеры репозиториев
│   │   └── queries/       # SQL-запросы (.sql)
│   └── mocks/             # сгенерированные моки (mockgen)
├── pkg/                   # переиспользуемые пакеты (планируются к выносу в отдельные репозитории)
│   ├── db/                # pgxpool, lifecycle, транзакционный менеджер
│   ├── health/            # health-checker и HTTP-хендлер
│   └── server/            # HTTP-сервер и middleware
├── migrations/
│   ├── db/postgres/       # миграции в формате Goose (управляются через Atlas)
│   └── current_schema.sql # актуальная схема БД (используется в тестах)
├── testhelpers/           # общие хелперы для unit и e2e тестов
└── tests/e2e/             # end-to-end тесты
```

## Слои и зависимости

Зависимости направлены строго внутрь. Внешние слои знают о внутренних, обратное запрещено.

```
handler → usecase → domain
                 ↑
          dbstorage (реализует интерфейсы domain)
```

**domain** — ядро системы. Содержит тип `Product`, правила валидации и интерфейс `ProductRepository`. Не зависит ни от чего внешнего — ни от HTTP, ни от БД, ни от сторонних библиотек кроме `lib-validation`.

**usecase** — оркестрирует доменные операции. Каждый use case реализует единый интерфейс:
```go
type Executer[C, R any] interface {
    Execute(ctx context.Context, cmd C) result.Value[R]
}
```
Use cases не знают об HTTP и не знают о конкретной реализации хранилища — только об интерфейсе `domain.ProductRepository`.

**handler** — транспортный слой. Декодирует HTTP-запрос, вызывает use case, формирует RFC 9457 ответ. Использует `lib-response` для сериализации и `lib-result` для обработки результата.

**dbstorage** — реализует `domain.ProductRepository` поверх sqlc-сгенерированного кода. Маппит pgx-ошибки в `apperr.Kind` (например, нарушение уникального индекса → `KindConflict`). Поддерживает транзакции через `ctxval` — tx-объект передаётся через контекст без изменения сигнатур методов.

## Обработка ошибок

Ошибки классифицируются через `lib-apperr`. Каждая ошибка несёт `Kind` — семантическую метку, по которой транспортный слой выбирает HTTP-статус.

```
domain / usecase           handler (lib-response)
─────────────────          ──────────────────────
KindNotFound       ──────► 404 Not Found
KindConflict       ──────► 409 Conflict
KindUnprocessable  ──────► 422 Unprocessable Entity
KindUnavailable    ──────► 503 Service Unavailable
KindInternal       ──────► 500 (detail скрывается)
*validation.Error  ──────► 422 + extensions.fields
```

Доменный слой не знает об HTTP-кодах. Транспортный слой не разбирает строки ошибок.

## Валидация

Валидация проходит на двух уровнях.

На уровне HTTP-запроса (`handler/validator.go`) проверяются структурные ограничения: обязательные поля, длина строк, допустимые значения enum. Используется `go-playground/validator` с маппингом ошибок в `validation.FieldError`.

На уровне домена (`domain/product.go`) проверяются бизнес-правила: формат slug через regexp, диапазон `price_cents`, формат `currency` по ISO 4217. Используется `lib-validation.Collector` для накопления всех ошибок за один проход.

Оба уровня возвращают `*validation.Error`, который `lib-response` сериализует в `extensions.fields`.

## Конфигурация

Конфигурация загружается через `lib-config` в два этапа: сначала применяются хардкодные значения по умолчанию из `cmd/ms-catalog/values.go`, затем переменные окружения с префиксом `INIT__` их перекрывают. Все ошибки конфигурации выдаются одновременно через `errors.Join` до старта сервера.

Двойное подчёркивание `__` используется как разделитель вложенности: `INIT__DB__DSN` → `db.dsn`.

## DI и жизненный цикл

Граф зависимостей собирается через `uber/fx`. Каждый компонент с состоянием (`*db.Db`, `*server.Server`) регистрирует `fx.StartStopHook` для корректного старта и graceful shutdown.

Порядок старта: конфигурация → логгер → пул БД (ping) → роутер → HTTP-сервер.  
Порядок остановки: HTTP-сервер (drain) → пул БД (close).

## Тестирование

**Unit-тесты** (`*_test.go` рядом с кодом) покрывают domain-валидацию, use cases и HTTP-обработчики. Use cases тестируются через `MockProductRepository` (mockgen). HTTP-обработчики тестируются через `MockExecuter` без поднятия сервера.

**Integration-тесты** (`internal/db/dbstorage/*_test.go`) проверяют адаптеры против реальной БД в testcontainer. Каждый тест получает собственную транзакцию через `dbtest.Suite`, которая откатывается в `TearDownTest` — изоляция без TRUNCATE.

**End-to-end тесты** (`tests/e2e/`) запускают полный стек через `httptest.Server` с реальной БД в testcontainer. Контейнер и пул разделяются между тестами через `sync.Once`.

## Инфраструктура

Образ собирается через [ko](https://ko.build) — без Dockerfile, из Go-бинаря в distroless-базу.

Миграции управляются через [Atlas](https://atlasgo.io) с форматом Goose. Схема описывается декларативно в `migrations/current_schema.sql`, Atlas вычисляет diff и генерирует миграцию.

В dev-окружении используется [air](https://github.com/air-verse/air) для hot-reload через Docker Compose.