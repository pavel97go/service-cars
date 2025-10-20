package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connString := "postgres://cars:cars@localhost:5432/cars?sslmode=disable"

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer pool.Close()

	sql := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS cars (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		brand VARCHAR(50) NOT NULL,
		model VARCHAR(50) NOT NULL,
		year INT NOT NULL CHECK (year >= 1886 AND year <= EXTRACT(YEAR FROM CURRENT_DATE) + 1),
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`
	_, err = pool.Exec(ctx, sql)
	if err != nil {
		log.Fatalf("migration error: %v", err)
	}
	fmt.Println("✅ migration applied successfully")
}
