package app

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/pavel97go/service-cars/internal/config"
	"github.com/pavel97go/service-cars/internal/handler"
	"github.com/pavel97go/service-cars/internal/metrics"
	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/repository"
	"github.com/pavel97go/service-cars/internal/router"
	"github.com/pavel97go/service-cars/internal/storage"
	"github.com/pavel97go/service-cars/internal/usecase"
)

func Run(ctx context.Context) error {
	models.Validate()
	metrics.Init()

	cfg := config.Init()
	addr := ":" + cfg.App.Port
	dsn := cfg.GetConnStr()

	pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := storage.GetConnect(pctx, dsn)
	if err != nil {
		return err
	}
	defer pool.Close()

	repo := repository.NewCarRepo(pool)
	uc := usecase.NewCarUsecase(repo)
	h := handler.NewCarHandler(uc)

	app := fiber.New()
	app.Use(metrics.Middleware())
	app.Get("/metrics", metrics.Handler())
	router.Register(app, h)

	log.Printf("Server is running on %s", addr)
	return app.Listen(addr)
}
