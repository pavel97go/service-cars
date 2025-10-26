# Service-Cars

Сервис для управления информацией об автомобилях.  
Реализована **чистая архитектура**, логирование, метрики Prometheus и трейсинг через Jaeger.

---

## Технологический стек
- **Golang (Fiber, PGX, Validator, Goose, Prometheus, OpenTelemetry)**
- **PostgreSQL + Docker Compose**
- **Jaeger** для распределённого трейсинга
- **Prometheus** для метрик (эндпоинт `/metrics`)
- **Makefile** для удобного запуска команд

---

## Архитектура проекта
```
cmd/                # Точка входа приложения
internal/
├── app/            # Инициализация зависимостей
├── handler/        # HTTP-обработчики (Fiber)
├── usecase/        # Бизнес-логика
├── repository/     # Работа с базой данных (PostgreSQL)
├── cache/          # In-memory кеш
├── metrics/        # Prometheus middleware
├── tracing/        # OpenTelemetry Jaeger
└── models/         # DTO и доменные модели
```

---

## Шаги по запуску проекта

### 1. Поднять контейнеры
```bash
make docker-up
```

### 2. Применить миграции
```bash
make migrate
```

### 3. Запустить приложение
```bash
go run ./cmd
```

После запуска будут доступны:
- API: http://localhost:8080/api/v1/cars
- Метрики приложения: http://localhost:8080/metrics  
  *(опционально можно подключить Prometheus и нацелить его на `/metrics`)*
- Jaeger UI: http://localhost:16686

---

## Команды, используемые в Makefile
| Команда | Описание |
|----------|-----------|
| `make run` | Запуск приложения |
| `make build` | Сборка бинаря |
| `make test` | Запуск тестов |
| `make lint` | Проверка линтером |
| `make migrate` | Применение миграций |
| `make migrate-add name=<name>` | Создание новой миграции |
| `make docker-up` | Запуск контейнеров |
| `make docker-down` | Остановка контейнеров |
| `make race` | Проверка на race conditions |

---

## REST API

| Метод | Эндпоинт | Описание |
|--------|-----------|-----------|
| `POST` | `/api/v1/cars/` | Создать автомобиль |
| `GET` | `/api/v1/cars/` | Получить список автомобилей |
| `GET` | `/api/v1/cars/:id` | Получить авто по ID |
| `PATCH` | `/api/v1/cars/:id` | Частично обновить данные автомобиля |
| `DELETE` | `/api/v1/cars/:id` | Удалить автомобиль |

Пример запроса:
```bash
curl -X POST http://localhost:8080/api/v1/cars \
  -H "Content-Type: application/json" \
  -d '{"brand":"Toyota","model":"Camry","year":2024}'
```

---

## Prometheus
Доступны по адресу:  
`http://localhost:8080/metrics`

Основные метрики:
- `http_requests_total` — количество HTTP-запросов
- `http_request_duration_seconds` — время ответа

---

## Тесты
```bash
make test
make race
```

Покрытие тестами (на текущий момент):
- **cache:** ~87%
- **usecase:** ~15%

---
