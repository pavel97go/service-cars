package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/pavel97go/service-cars/internal/apperr"
	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/repository"
)

type CarUC struct {
	repo repository.CarProvider
}

func NewCarUsecase(repo repository.CarProvider) CarUsecase {
	return &CarUC{repo: repo}
}
func (u *CarUC) Create(ctx context.Context, req models.CreateCarRequest) (models.CarResponse, error) {
	if err := models.ValidateStruct(req); err != nil {
		return models.CarResponse{}, err
	}
	currentYear := time.Now().Year()
	if req.Year > currentYear+1 {
		return models.CarResponse{}, fmt.Errorf("%w: year %d invalid", apperr.ErrInvalidInput, req.Year)
	}
	car := models.Car{
		Brand: req.Brand,
		Model: req.Model,
		Year:  req.Year,
	}
	if err := u.repo.InsertCar(ctx, &car); err != nil {
		return models.CarResponse{}, err
	}
	return models.CarResponse{
		ID:    car.ID,
		Brand: car.Brand,
		Model: car.Model,
		Year:  car.Year,
	}, nil
}
func (u *CarUC) List(ctx context.Context) ([]models.CarResponse, error) {
	cars, err := u.repo.ListCars(ctx)
	if err != nil {
		return nil, err
	}
	return models.ToResponse(cars), nil
}
func (u *CarUC) Get(ctx context.Context, id string) (models.CarResponse, error) {
	car, err := u.repo.GetCarByID(ctx, id)
	if err == apperr.ErrNotFound {
		return models.CarResponse{}, apperr.ErrNotFound
	}
	if err != nil {
		return models.CarResponse{}, err
	}
	response := models.CarResponse{
		ID:    car.ID,
		Brand: car.Brand,
		Model: car.Model,
		Year:  car.Year,
	}
	return response, nil
}
func (u *CarUC) Update(ctx context.Context, req models.UpdateCarRequest) (models.CarResponse, error) {
	if err := models.ValidateStruct(req); err != nil {
		return models.CarResponse{}, err
	}
	car, err := u.repo.GetCarByID(ctx, req.ID)
	if err == apperr.ErrNotFound {
		return models.CarResponse{}, apperr.ErrNotFound
	}
	if req.Brand != "" {
		car.Brand = req.Brand
	}
	if req.Model != "" {
		car.Model = req.Model
	}
	if req.Year != 0 {
		yearLimit := time.Now().Year() + 1
		if req.Year > yearLimit {
			return models.CarResponse{}, fmt.Errorf("%w: year must be <= %d", apperr.ErrInvalidInput, yearLimit)
		}
		car.Year = req.Year
	}

	if err := u.repo.UpdateCar(ctx, car); err != nil {
		if err == apperr.ErrNotFound {
			return models.CarResponse{}, apperr.ErrNotFound
		}
		return models.CarResponse{}, err
	}
	return models.CarResponse{
		ID:    car.ID,
		Brand: car.Brand,
		Model: car.Model,
		Year:  car.Year,
	}, nil
}
func (u *CarUC) Delete(ctx context.Context, id string) error {
	err := u.repo.DeleteByID(ctx, id)
	if err == apperr.ErrNotFound {
		return apperr.ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
