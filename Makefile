ifneq (,$(wildcard ./.env))
   include ./.env
   TO_EXPORT := $(filter KO_%,POSTGRES_%,$(.VARIABLES))
   export $(TO_EXPORT)
endif

.PHONY: help
help:
	@awk -f help.awk $(MAKEFILE_LIST)
	@echo ""

.DEFAULT_GOAL := help

APP_NAME         := ms-catalog
CMD_PATH         := ./cmd/$(APP_NAME)
CMD_COMPOSE_PROD := docker compose --project-directory . -f deploy/compose.yml
CMD_COMPOSE_DEV  := docker compose --project-directory . -f deploy/compose.hot-reload.yml

VERSION    ?= $(shell git describe --tags --exact-match 2>/dev/null || echo dev)
COMHASH    ?= $(shell git rev-parse HEAD 2>/dev/null || echo "undefined")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || echo "undefined")
GO_VERSION ?= $(shell go mod edit -json | jq -r '.Go' | cut -d'.' -f1,2)

ATLAS_MIGRATION_DIR=file://migrations/db/postgres?format=goose
ATLAS_DEV_URL=docker://postgres?search_path=public
ATLAS_SCHEMA_FILE=file://migrations/current_schema.sql

DB_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5433/$(POSTGRES_DB)?sslmode=disable

# === Production ===

.PHONY: prod
prod: prod-build prod-up ## Собрать и запустить prod окружение

.PHONY: prod-build
prod-build: ## Собрать prod образ через ko
	ko build $(CMD_PATH) -B -L

.PHONY: prod-up
prod-up: ## Запустить prod окружение
	$(CMD_COMPOSE_PROD) up -d

.PHONY: prod-down
prod-down: ## Остановить prod окружение
	$(CMD_COMPOSE_PROD) down

.PHONY: prod-logs
prod-logs: ## Показать логи prod окружения
	$(CMD_COMPOSE_PROD) logs -f

.PHONY: prod-restart
prod-restart: prod-down prod ## Перезапустить prod окружение

# === Development (hot-reload) ===

.PHONY: dev
dev: dev-build dev-up ## Собрать и запустить dev окружение

.PHONY: dev-build
dev-build: ## Пересобрать dev образ
	$(CMD_COMPOSE_DEV) build --build-arg GO_VERSION=$(or $(GO_VERSION),1.26)

.PHONY: dev-up
dev-up: ## Запустить dev окружение (hot-reload)
	$(CMD_COMPOSE_DEV) up -d

.PHONY: dev-down
dev-down: ## Остановить dev окружение
	$(CMD_COMPOSE_DEV) down

.PHONY: dev-logs
dev-logs: ## Показать логи dev окружения
	$(CMD_COMPOSE_DEV) logs -f

.PHONY: dev-restart
dev-restart: dev-down dev ## Перезапустить dev окружение

# === Generation ===

.PHONY: code-gen
code-gen: ## Сгенерировать код
	go generate ./...

.PHONY: sqlc-gen
sqlc-gen: ## Сгенерировать Go-код из queries
	sqlc generate --file .sqlc.yml

.PHONY: swag-gen
swag-gen: ## Сгенерировать OpenAPI spec
	swag init \
		--generalInfo main.go \
		--dir         $(CMD_PATH),internal/handler \
		--output      docs \
		--outputTypes go,json,yaml \
		--parseDependency \
		--parseInternal

# === CI/CD ===

.PHONY: lint
lint: ## Запустить линтер
	golangci-lint cache clean
	golangci-lint run

.PHONY: test
test: ## Запустить тесты с покрытием
	go generate ./...
	go test -race -coverprofile=coverage.out -covermode=atomic \
		$(shell go list ./internal/... | grep -v -E '/(vendor|gen|cmd|testdata|mocks)(/|$$)')
	sed -i -E '/\/gen\/|_mock\.go:|\.pb(\.gw)?\.go:|\.sql\.go:|\.gen\.go:/d' coverage.out || true
	go tool cover -func=coverage.out | tee coverage.humanize | tail -1

.PHONY: test-e2e
test-e2e: ## Запустить end-to-end тесты
	go test -v -timeout 120s ./tests/e2e/...

# === Profiling ===

.PHONY: prof
prof: ## Собрать профили (cpu, mem, block, mutex)
	go test -run=. -bench=. \
		-cpuprofile=cpu.out -blockprofile=block.out \
		-memprofile=mem.out -mutexprofile=mutex.out ./...

.PHONY: prof-view
prof-view: ## Открыть профиль в браузере (e.g. make prof-view FILE=mem.out)
	go tool pprof -http=:8082 $(or $(FILE),cpu.out)

# === Migrations ===

.PHONY: mig-diff
mig-diff: ## Cоздать новую миграцию (использование: make mig-diff NAME=add_user_table)
	@if [ -z "$(NAME)" ]; then echo "Usage: make mig-diff NAME=add_user_table"; exit 1; fi
	@atlas migrate diff $(NAME) --dir $(ATLAS_MIGRATION_DIR) --to $(ATLAS_SCHEMA_FILE) --dev-url $(ATLAS_DEV_URL)
	@atlas migrate hash --force --dir $(ATLAS_MIGRATION_DIR)

.PHONY: mig-lint
mig-lint: ## Проверка качества миграций
	@atlas migrate lint --dir "$(ATLAS_MIGRATION_DIR)" --dev-url "$(ATLAS_DEV_URL)" --latest 1

.PHONY: mig-apply
mig-apply: ## Применить миграции
	@atlas migrate apply --dir "$(ATLAS_MIGRATION_DIR)" --url "$(DB_DSN)"

.PHONY: mig-valid
mig-valid: ## Валидация миграций
	@atlas migrate validate --dir "$(ATLAS_MIGRATION_DIR)" --dev-url "$(ATLAS_DEV_URL)"

.PHONY: mig-status
mig-status: ## Показать статус миграций
	@atlas migrate status --dir "$(ATLAS_MIGRATION_DIR)" --url "$(DB_DSN)"

# === Utils ===

.PHONY: clean
clean: ## Очистить артефакты
	rm -rf bin/ tmp/ sbom/ coverage.*
	$(CMD_COMPOSE_PROD) down -v 2>/dev/null || true
	$(CMD_COMPOSE_DEV)  down -v 2>/dev/null || true

.PHONY: down
down: ## Остановить все окружения
	$(CMD_COMPOSE_PROD) down 2>/dev/null || true
	$(CMD_COMPOSE_DEV)  down 2>/dev/null || true

.PHONY: version
version: ## Показать версию
	@echo "Version:    $(VERSION)"
	@echo "Comhash:    $(COMHASH)"
	@echo "Build Date: $(BUILD_DATE)"
	@echo "Go Version: $(GO_VERSION)"