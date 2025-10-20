package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pavel97go/service-cars/internal/handler"
	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/repository"
	"github.com/pavel97go/service-cars/internal/router"
	"github.com/pavel97go/service-cars/internal/usecase"
)

func Run(ctx context.Context) error {
	models.Validate()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://cars:cars@localhost:5432/cars?sslmode=disable"
	}
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(pctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()
	repo := repository.NewCarRepo(pool)
	uc := usecase.NewCarUsecase(repo)
	h := handler.NewCarHandler(uc)

	app := fiber.New()
	router.Register(app, h)

	log.Printf(" Server running on %s", addr)
	return app.Listen(addr)
}
