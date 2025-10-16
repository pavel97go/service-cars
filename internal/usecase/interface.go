package usecase

import (
	"context"

	"github.com/pavel97go/service-cars/internal/models"
)

type CarUsecase interface {
	Create(ctx context.Context, req models.CreateCarRequest) (models.CarResponse, error)
	List(ctx context.Context) ([]models.CarResponse, error)
	Get(ctx context.Context, id string) (models.CarResponse, error)
	Update(ctx context.Context, req models.UpdateCarRequest) (models.CarResponse, error)
	Delete(ctx context.Context, id string) error
}
