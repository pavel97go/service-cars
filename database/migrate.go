package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var migrations embed.FS

func Migrate(connStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return errors.Wrap(err, "sql.Open(pgx)")
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("db.Close error:", err)
		}
	}()

	if err := db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping db")
	}
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return errors.Wrap(err, "set dialect")
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return errors.Wrap(err, "apply migrations")
	}
	return nil
}
