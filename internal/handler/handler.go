package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/pavel97go/service-cars/internal/apperr"
	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/usecase"
)

type CarHandler struct {
	uc usecase.CarUsecase
}

func NewCarHandler(uc usecase.CarUsecase) *CarHandler {
	return &CarHandler{uc: uc}
}

func (h *CarHandler) Create(c *fiber.Ctx) error {
	var req models.CreateCarRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.uc.Create(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	c.Location("/api/v1/cars/" + resp.ID)
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *CarHandler) List(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.uc.List(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CarHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format, must be UUID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.uc.Get(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "car not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CarHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format, must be UUID"})
	}

	var req models.UpdateCarRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	req.ID = id

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.uc.Update(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, apperr.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "car not found"})
		case errors.Is(err, apperr.ErrInvalidInput):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *CarHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing id"})
	}
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id format, must be UUID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.uc.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, apperr.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "car not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}
	return c.SendStatus(fiber.StatusNoContent)
}
