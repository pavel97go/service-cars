package repository

import (
	"context"

	"github.com/pavel97go/service-cars/internal/models"
)

type CarProvider interface {
	ListCars(ctx context.Context) ([]models.Car, error)
	GetCarByID(ctx context.Context, id string) (*models.Car, error)
	InsertCar(ctx context.Context, newCar *models.Car) error
	UpdateCar(ctx context.Context, updatedCar *models.Car) error
	DeleteByID(ctx context.Context, id string) error
}
