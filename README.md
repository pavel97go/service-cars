# 🚗 Service Cars

**Service Cars** — это учебно-практический backend-проект на **Go (Golang)**, разработанный в стиле **Clean Architecture**.  
Сервис моделирует систему управления автопарком — CRUD-операции с машинами, in-memory кеширование, метрики и трейсинг в production-формате.

---

## Ключевая идея

Проект построен вокруг принципов чистой архитектуры: каждый слой отвечает строго за своё —  
от репозитория до бизнес-логики и хендлеров.  
Добавлены **Prometheus-метрики**, **OpenTelemetry-трейсинг** и **in-memory TTL-кеш**,  
чтобы показать владение продакшн-практиками.

---

## Технологический стек

| Категория | Технологии |
|------------|-------------|
| Язык | [Go 1.22+](https://go.dev/) |
| Web-фреймворк | [Fiber v2](https://gofiber.io/) |
| База данных | PostgreSQL + [pgxpool](https://github.com/jackc/pgx) |
| Миграции | [goose](https://github.com/pressly/goose) |
| Кеш | In-memory TTL Cache |
| Метрики | Prometheus |
| Трейсинг | OpenTelemetry + Jaeger |
| Тестирование | GoMock + Testify |
| Контейнеризация | Docker, docker-compose |

---

## Архитектура проекта

```
internal/
├── app/          # Точка входа и сборка зависимостей
├── config/       # ENV и конфигурация
├── handler/      # HTTP-обработчики (Fiber)
├── usecase/      # Бизнес-логика
├── repository/   # Работа с БД (Postgres)
├── cache/        # In-memory кеш с TTL
├── metrics/      # Prometheus метрики
├── tracing/      # OpenTelemetry трейсинг
├── models/       # DTO + валидация (validator/v10)
└── router/       # Регистрация маршрутов
```

---

## Запуск

### 🔹 Через Docker
```bash
docker-compose up --build
```
После запуска:
- API: [http://localhost:8080](http://localhost:8080)
- Prometheus: [http://localhost:9090](http://localhost:9090)
- Jaeger UI: [http://localhost:16686](http://localhost:16686)
- Метрики: [http://localhost:8080/metrics](http://localhost:8080/metrics)

---

### 🔹 Без Docker
```bash
go run ./cmd
```
Перед запуском — проверь `.env` или `internal/config/config.yml`!!

---

## API эндпоинты(entry-point)

| Метод | Путь | Описание |
|--------|------|----------|
| `GET`    | `/api/v1/cars`       | Получить список всех машин |
| `GET`    | `/api/v1/cars/:id`   | Получить машину по ID |
| `POST`   | `/api/v1/cars`       | Создать новую машину |
| `PUT`    | `/api/v1/cars/:id`   | Обновить данные машины |
| `DELETE` | `/api/v1/cars/:id`   | Удалить машину |

---

## Метрики и трейсинг

- `/metrics` — метрики в формате Prometheus  
- Jaeger — визуализация трейсинга запросов  
- OpenTelemetry — автоматический сбор спанов  

---

## Тестирование

```bash
go test ./...
```

Тесты включают мок-реализацию репозитория (GoMock) и проверку кеш-поведения.

---

## Makefile команды

```bash
make build          # Собрать проект
make run            # Запуск локально
make migrate-up     # Применить миграции
make migrate-add name=create_table
```

---

## Основные фишки проекта

*Чистая архитектура и разделение слоёв  
*Конфиг и окружение через ENV  
*In-memory кеш с TTL и конкурентной безопасностью  
*Prometheus-метрики и OpenTelemetry трейсинг  
*Мок-тесты для usecase и кеш-логики   
*Production-ready структура проекта  

---

## Автор проекта

**Pavel (lordiiikxd)**  
Backend Engineer | Golang Developer  
[github.com/pavel97go](https://github.com/pavel97go)  

---

⭐️ Если тебе зашёл проект — поставь звездочку!
