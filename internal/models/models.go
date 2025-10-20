package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func Validate() {
	validate = validator.New()
}

func ValidateStruct(v interface{}) error {
	return validate.Struct(v)
}

type CreateCarRequest struct {
	Brand string `json:"brand" validate:"required,alphaunicode,min=1,max=50"`
	Model string `json:"model" validate:"required,alphaunicode,min=1,max=50"`
	Year  int    `json:"year" validate:"required,gte=1886"`
}

type UpdateCarRequest struct {
	ID    string `json:"id" validate:"required,uuid4"`
	Brand string `json:"brand" validate:"omitempty,min=1,max=50"`
	Model string `json:"model" validate:"omitempty,min=1,max=50"`
	Year  int    `json:"year" validate:"omitempty,gte=1886"`
}

type Car struct {
	ID        string    `db:"id"`
	Brand     string    `db:"brand"`
	Model     string    `db:"model"`
	Year      int       `db:"year"`
	CreatedAt time.Time `db:"created_at"`
}

type CarResponse struct {
	ID    string `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	Year  int    `json:"year"`
}

func ToResponse(dbCars []Car) []CarResponse {
	response := make([]CarResponse, len(dbCars))
	for i, car := range dbCars {
		response[i] = CarResponse{
			ID:    car.ID,
			Brand: car.Brand,
			Model: car.Model,
			Year:  car.Year,
		}
	}
	return response
}

func UpdatedCarDTO(updatedCars UpdateCarRequest) Car {
	updated := Car{}
	updated.ID = updatedCars.ID
	updated.Brand = updatedCars.Brand
	updated.Model = updatedCars.Model
	updated.Year = updatedCars.Year

	return updated
}
