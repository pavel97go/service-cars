APP_NAME=service-cars
MIGRATIONS_DIR=database/migrations
DSN="postgres://cars:cars@localhost:5432/cars?sslmode=disable"

.PHONY: run build test migrate migrate-add docker-up docker-down

run: 
	go run ./cmd

build: 
	go build -o bin/$(APP_NAME) ./cmd

test: 
	go test ./... -cover

migrate: 
	goose -dir $(MIGRATIONS_DIR) pgx $(DSN) up

migrate-add: 
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

docker-up: 
	docker compose up -d

docker-down: 
	docker compose down