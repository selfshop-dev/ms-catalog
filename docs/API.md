# API Reference

Базовый путь: `/api/v1`

Все ответы используют единый JSON-конверт в формате [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457).
Успешные ответы несут поле `data`, ошибочные — `detail` и опционально `extensions.fields`.

## Формат ответов

**Успешный ответ (201 Created):**
```json
{
    "type":     "about:blank",
    "status":   201,
    "title":    "Created",
    "data":     { "id": "01960000-0000-7000-0000-000000000000" },
    "instance": "/api/v1/products",
    "meta":     { "request_id": "abc123" }
}
```

**Ошибочный ответ (422 Unprocessable Entity):**
```json
{
    "type":       "about:blank",
    "status":     422,
    "title":      "Unprocessable Entity",
    "detail":     "request",
    "instance":   "/api/v1/products",
    "meta":       { "request_id": "abc123" },
    "extensions": {
        "fields": [
            { "field": "name", "code": "required", "message": "name is required" },
            { "field": "slug", "code": "required", "message": "slug is required" }
        ]
    }
}
```

`Content-Type` для успешных ответов — `application/json`, для ошибочных — `application/problem+json`.

Каждый ответ содержит заголовок `X-Request-Id` для трассировки запроса.

## Товары

### POST /products

Создать новый товар.

**Тело запроса:**

| Поле | Тип | Обязательно | Описание |
|---|---|---|---|
| `name` | string | да | Название товара. От 1 до 128 символов |
| `slug` | string | да | URL-идентификатор. От 1 до 128 символов, только строчные буквы, цифры и дефисы |
| `description` | string | нет | Полное описание |
| `short_description` | string | нет | Краткое описание. До 256 символов |
| `display_image_url` | string | нет | URL изображения. Должен начинаться с `http://` или `https://` |
| `price_cents` | integer | нет | Цена в минимальных единицах валюты. Неотрицательное число, по умолчанию `0` |
| `currency` | string | нет | Код валюты ISO 4217. Ровно 3 заглавные буквы, по умолчанию `USD` |
| `status` | string | нет | Статус товара. Одно из: `active`, `inactive`, `draft`, `archived`. По умолчанию `active` |

**Пример запроса:**
```json
{
    "name":        "Wireless Keyboard",
    "slug":        "wireless-keyboard",
    "description": "Full-size wireless keyboard with USB receiver",
    "price_cents": 4999,
    "currency":    "USD",
    "status":      "draft"
}
```

**Ответы:**

| Статус | Описание |
|---|---|
| `201 Created` | Товар создан. Тело содержит `{ "id": "<uuid>" }` |
| `400 Bad Request` | Некорректное тело запроса (невалидный JSON) |
| `422 Unprocessable Entity` | Ошибка валидации полей. Тело содержит `extensions.fields` |
| `409 Conflict` | Товар с таким slug уже существует |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

---

### GET /products/active

Получить постраничный список активных товаров, отсортированных по дате создания (новые первые).

**Query-параметры:**

| Параметр | Тип | По умолчанию | Описание |
|---|---|---|---|
| `limit` | integer | `20` | Количество записей на странице. От 1 до 100 |
| `offset` | integer | `0` | Смещение от начала списка |

**Пример:**
```
GET /api/v1/products/active?limit=10&offset=20
```

**Ответы:**

| Статус | Описание |
|---|---|
| `200 OK` | Список товаров. Тело содержит `items`, `limit`, `offset` |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

**Пример тела ответа:**
```json
{
    "items": [
        {
            "id":          "01960000-0000-7000-0000-000000000000",
            "name":        "Wireless Keyboard",
            "slug":        "wireless-keyboard",
            "price_cents": 4999,
            "currency":    "USD",
            "status":      "active",
            "created_at":  "2026-04-12T10:00:00Z",
            "updated_at":  "2026-04-12T10:00:00Z"
        }
    ],
    "limit":  10,
    "offset": 20
}
```

---

### GET /products/{id}

Получить товар по UUID.

**Path-параметры:**

| Параметр | Тип | Описание |
|---|---|---|
| `id` | uuid | UUID товара |

**Ответы:**

| Статус | Описание |
|---|---|
| `200 OK` | Товар найден |
| `400 Bad Request` | Некорректный формат UUID |
| `404 Not Found` | Товар не найден |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

---

### GET /products/{slug}/by-slug

Получить товар по slug.

**Path-параметры:**

| Параметр | Тип | Описание |
|---|---|---|
| `slug` | string | Slug товара |

**Ответы:**

| Статус | Описание |
|---|---|
| `200 OK` | Товар найден |
| `400 Bad Request` | Пустой slug |
| `404 Not Found` | Товар не найден |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

**Пример тела ответа (200):**
```json
{
    "id":                "01960000-0000-7000-0000-000000000000",
    "name":              "Wireless Keyboard",
    "slug":              "wireless-keyboard",
    "description":       "Full-size wireless keyboard with USB receiver",
    "short_description": null,
    "display_image_url": null,
    "price_cents":       4999,
    "currency":          "USD",
    "status":            "active",
    "created_at":        "2026-04-12T10:00:00Z",
    "updated_at":        "2026-04-12T10:00:00Z"
}
```

---

### PUT /products/{id}

Полностью заменить поля товара. Статус не изменяется — для этого используется отдельный эндпоинт.

**Path-параметры:**

| Параметр | Тип | Описание |
|---|---|---|
| `id` | uuid | UUID товара |

**Тело запроса:**

| Поле | Тип | Обязательно | Описание |
|---|---|---|---|
| `name` | string | да | Название. От 1 до 128 символов |
| `slug` | string | да | Slug. От 1 до 128 символов |
| `currency` | string | да | Код валюты ISO 4217 |
| `description` | string | нет | Полное описание |
| `short_description` | string | нет | Краткое описание. До 256 символов |
| `display_image_url` | string | нет | URL изображения |
| `price_cents` | integer | нет | Цена в минимальных единицах валюты |

**Ответы:**

| Статус | Описание |
|---|---|
| `200 OK` | Товар обновлён. Тело содержит `{ "id": "<uuid>" }` |
| `400 Bad Request` | Некорректный UUID или тело запроса |
| `404 Not Found` | Товар не найден |
| `409 Conflict` | Товар с таким slug уже существует |
| `422 Unprocessable Entity` | Ошибка валидации полей |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

---

### PUT /products/{id}/status

Изменить статус товара.

**Path-параметры:**

| Параметр | Тип | Описание |
|---|---|---|
| `id` | uuid | UUID товара |

**Тело запроса:**

| Поле | Тип | Обязательно | Описание |
|---|---|---|---|
| `status` | string | да | Новый статус. Одно из: `active`, `inactive`, `draft`, `archived` |

**Пример:**
```json
{ "status": "active" }
```

**Ответы:**

| Статус | Описание |
|---|---|
| `200 OK` | Статус изменён. Тело содержит `{ "id": "<uuid>" }` |
| `400 Bad Request` | Некорректный UUID или тело запроса |
| `404 Not Found` | Товар не найден |
| `422 Unprocessable Entity` | Недопустимое значение статуса |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

---

### DELETE /products/{id}

Мягко удалить товар (устанавливает `deleted_at`). Удалённый товар не возвращается ни одним GET-эндпоинтом.

**Path-параметры:**

| Параметр | Тип | Описание |
|---|---|---|
| `id` | uuid | UUID товара |

**Ответы:**

| Статус | Описание |
|---|---|
| `204 No Content` | Товар удалён |
| `400 Bad Request` | Некорректный формат UUID |
| `404 Not Found` | Товар не найден |
| `500 Internal Server Error` | Внутренняя ошибка сервера |

---

## Служебные эндпоинты

### GET /health/alive

Проверка жизнеспособности процесса. Возвращает `200 OK` если HTTP-сервер принимает соединения.

```json
{ "status": "ok" }
```

### GET /health/ready

Проверка готовности сервиса к приёму трафика. Запускает все зарегистрированные health-чекеры (PostgreSQL).

```json
{
    "status": "ok",
    "checkers": [
        { "name": "db::postgres" }
    ]
}
```

При сбое любого чекера возвращает `503 Service Unavailable`:

```json
{
    "status": "unhealthy",
    "checkers": [
        { "name": "db::postgres", "error": "dial tcp: connection refused" }
    ]
}
```

### GET /log/level

Получить текущий уровень логирования.

```json
{ "level": "info" }
```

### PUT /log/level

Изменить уровень логирования без перезапуска сервиса.

```json
{ "level": "debug" }
```

Допустимые уровни: `debug`, `info`, `warn`, `error`, `panic`, `fatal`.

### GET /swagger

Swagger UI с интерактивной документацией API. Доступен только в `runmode=dev`.

### GET /swagger/doc.json

OpenAPI 2.0 спецификация в формате JSON. Доступна только в `runmode=dev`.