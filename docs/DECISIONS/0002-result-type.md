# 0002 — Result[T, E] в use cases

### Контекст

Use cases принимают команду и возвращают результат или ошибку. Стандартная Go-пара `(T, error)` в сигнатуре интерфейса `Executer[C, R]` приводит к тому, что вызывающий код обязан немедленно обработать ошибку — нельзя передать результат дальше без распаковки.

В handler-слое это выражается в повторяющемся паттерне:

```go
res, err := u.Execute(r.Context(), cmd)
if err != nil {
    handler.Respond.Error(w, r, err)
    return
}
// использовать res
```

## Решение

Use cases возвращают `result.Value[R]` — алиас для `Result[R, error]` из `lib-result`. Интерфейс:

```go
type Executer[C, R any] interface {
    Execute(ctx context.Context, cmd C) result.Value[R]
}
```

Внутри usecase трансформации через `result.Map` и `result.Of` позволяют строить pipeline без явных `if err != nil`:

```go
func (u *getProduct) Execute(ctx context.Context, qry GetProductQuery) result.Value[GetProductResult] {
    return result.Map(
        result.Of(u.r.GetByID(ctx, qry.ID)),
        func(p *domain.Product) GetProductResult { return GetProductResult{Product: p} },
    )
}
```

На границе с handler-слоем результат разворачивается через `.ToGo()`:

```go
res, err := u.Execute(r.Context(), cmd).ToGo()
```

### Последствия

- Реализации use cases лаконичнее за счёт pipeline-стиля.
- Интерфейс `Executer` однороден — один метод, одно возвращаемое значение.
- В mock-тестах возвращаемое значение оборачивается в `result.Ok` или `result.Of`.

## Ограничения

`Result` добавляет нетипичный для Go стиль. Разработчики знакомые только со стандартным Go-идиомом могут воспринять его как усложнение. Использование ограничено слоем usecase — domain и dbstorage используют стандартную пару `(T, error)`.

Связанные решения: [0001 — Clean Architecture](0001-clean-architecture.md).