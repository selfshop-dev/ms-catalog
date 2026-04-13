# Contributing

Проект следует [`Conventional Commits`](https://www.conventionalcommits.org). Формат: `<type>: <summary>`, где `summary` — повелительное наклонение,английский язык, без точки в конце.

| Тип | Когда использовать |
|---|---|
| `feat` | Новая функциональность |
| `fix` | Исправление бага |
| `refactor` | Рефакторинг без изменения поведения |
| `perf` | Оптимизация производительности |
| `test` | Добавление или изменение тестов |
| `docs` | Документация |
| `ci` | Изменения CI/CD и GitHub Actions |
| `build` | Сборка, зависимости, инфраструктура проекта |
| `chore` | Рутинные задачи, не попадающие в другие типы |

Тип коммита автоматически определяет метку PR и тип следующего релиза:  

- `feat!` / `fix!` (breaking change) → major
- `feat` → minor
- `fix` / `perf` / `refactor` → patch


## Makefile
 
Все команды запускаются из корня репозитория через `make <цель>`. Полный список с описаниями — `make help`.
 
### Production
 
| Цель | Описание |
|---|---|
| `make prod` | Собрать образ через ko и запустить prod окружение |
| `make prod-build` | Собрать prod образ через ko |
| `make prod-up` | Запустить prod окружение |
| `make prod-down` | Остановить prod окружение |
| `make prod-logs` | Показать логи prod окружения |
| `make prod-restart` | Перезапустить prod окружение |
 
### Development (hot-reload)
 
| Цель | Описание |
|---|---|
| `make dev` | Собрать образ и запустить dev окружение с hot-reload |
| `make dev-build` | Пересобрать dev образ |
| `make dev-up` | Запустить dev окружение |
| `make dev-down` | Остановить dev окружение |
| `make dev-logs` | Показать логи dev окружения |
| `make dev-restart` | Перезапустить dev окружение |
 
### Генерация кода
 
| Цель | Описание |
|---|---|
| `make code-gen` | Запустить `go generate ./...` |
| `make sqlc-gen` | Сгенерировать Go-код из SQL-запросов через sqlc |
| `make swag-gen` | Сгенерировать OpenAPI spec через swag |
 
### Тестирование и анализ
 
| Цель | Описание |
|---|---|
| `make test` | Генерация кода + unit-тесты с coverage |
| `make test-e2e` | End-to-end тесты (требует запущенной БД) |
| `make lint` | Запустить golangci-lint |
| `make prof` | Собрать профили (cpu, mem, block, mutex) |
| `make prof-view` | Открыть профиль в браузере (`FILE=mem.out`, по умолчанию `cpu.out`) |
 
### Миграции
 
| Цель | Описание |
|---|---|
| `make mig-diff NAME=<name>` | Создать новую миграцию по diff схемы |
| `make mig-apply` | Применить миграции к БД |
| `make mig-lint` | Проверить качество последней миграции |
| `make mig-valid` | Валидировать все миграции |
| `make mig-status` | Показать статус применённых миграций |
 
### Утилиты
 
| Цель | Описание |
|---|---|
| `make version` | Показать версию, commit hash, дату сборки |
| `make clean` | Удалить артефакты сборки и остановить контейнеры |
| `make down` | Остановить все окружения |

## Требования к коду

Подробно описаны в [`PULL_REQUEST_TEMPLATE`](https://github.com/selfshop-dev/.github/blob/main/PULL_REQUEST_TEMPLATE.md?plain=1) (раздел Checklist).

## Issues и уязвимости

- Общие вопросы и баги — через [`шаблоны issues`](../../issues/new/choose).
- **Уязвимости** — **только** через [`🔒 Private Vulnerability Reporting`](../../security/advisories/new). Не создавай публичный issue или discussion!