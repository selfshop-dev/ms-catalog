# 0003 — Транспорт-агностичная классификация ошибок

### Контекст

HTTP-обработчик должен выбрать статус ответа на основе ошибки из domain или usecase. Без явной классификации единственный способ — проверять тип ошибки или разбирать строку сообщения:

```go
if strings.Contains(err.Error(), "not found") { // хрупко
    http.Error(w, "", 404)
}
```

Это связывает транспортный слой с деталями реализации нижележащих слоёв и делает логику маппинга непредсказуемой.

## Решение

Все ошибки классифицируются через `lib-apperr` с помощью типа `Kind`. Domain и storage возвращают `*apperr.Error` с семантической меткой — без HTTP-зависимостей:

```go
// dbstorage — не знает об HTTP
if errors.Is(err, pgx.ErrNoRows) {
    return nil, apperr.NotFoundf("product not found with ID %s", id)
}
```

Handler-слой через `lib-response` маппит `Kind` в HTTP-статус:

```
KindNotFound       → 404
KindConflict       → 409
KindUnprocessable  → 422
KindUnavailable    → 503
KindInternal       → 500 (detail скрывается от клиента)
```

Ошибки валидации (`*validation.Error`) маппятся в 422 с полем `extensions.fields`.

### Последствия

- Domain и storage не зависят от `net/http`.
- Маппинг ошибок сосредоточен в одном месте — `lib-response`.
- `KindInternal` никогда не раскрывает детали клиенту.
- `KindUnavailable` и `KindTimeout` помечаются как retryable — клиент может принять решение о повторной попытке.

## Ограничения

Требует дисциплины: все ошибки из dbstorage должны быть явно классифицированы. Неклассифицированная ошибка попадает в `KindUnknown` → 500, что может скрывать реальные проблемы.

Связанные решения: [0001 — Clean Architecture](0001-clean-architecture.md), [0004 — RFC 9457](0004-rfc9457-problem-details.md).