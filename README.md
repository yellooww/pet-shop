# Go Pet Shop

Учебный проект интернет-магазина товаров для питомцев, написанный на Go с использованием PostgreSQL.

## Запуск приложения

Для запуска приложения настрой `.env` по примеру `.env.example`, примени миграции и запусти сервер:

```bash
task migrate
go run ./cmd/app
```

### Запуск тестов

Из корневой папки проекта:

```bash
go test ./...
```

## Две версии практики: v1 и v2

Практика хранится в двух Git-ветках.

### v1 — ручной мок

Ветка `v1` использует `internal/handlers/product/product_mock.go`. Вручную написанный `ProductsMock` реализует интерфейс `Products`.

### v2 — сгенерированный мок (mockery)

В `internal/handlers/product/products.go` находится директива:

```go
//go:generate go run github.com/vektra/mockery/v2 --name=Products
```

Она запускает `mockery`, который генерирует мок. Ветка `v2` использует сгенерированный файл в `internal/handlers/product/mocks/Products.go`; ручной `product_mock.go` в этой ветке не используется.

Перегенерация:

```bash
go generate ./...
```

В `v2` тесты импортируют `go-pet-shop/internal/handlers/product/mocks`. Логика проверок HTTP-статуса и тела ответа остаётся такой же, как в `v1`.
