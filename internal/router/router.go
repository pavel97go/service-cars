package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pavel97go/service-cars/internal/handler"
)

func Register(app *fiber.App, h *handler.CarHandler) {
	api := app.Group("api/v1")
	cars := api.Group("/cars")

	cars.Post("/", h.Create)
	cars.Get("/", h.List)
	cars.Get("/:id", h.Get)
	cars.Patch("/:id", h.Update)
	cars.Delete("/:id", h.Delete)
}
