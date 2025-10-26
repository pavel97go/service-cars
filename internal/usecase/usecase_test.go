package usecase_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/repository/mocks"
	"github.com/pavel97go/service-cars/internal/usecase"
)

func TestMain(m *testing.M) {
	models.Validate()
	os.Exit(m.Run())
}
func TestCreateCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockCarProvider(ctrl)
	uc := usecase.NewCarUsecase(mockRepo)

	req := models.CreateCarRequest{
		Brand: "Toyota",
		Model: "Camry",
		Year:  time.Now().Year(),
	}

	mockRepo.
		EXPECT().
		InsertCar(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, c *models.Car) error {
			c.ID = "uuid-1"
			c.CreatedAt = time.Now()
			return nil
		}).
		Times(1)

	resp, err := uc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Brand != "Toyota" {
		t.Fatalf("unexpected brand: got %q, want %q", resp.Brand, "Toyota")
	}

	if resp.Model != "Camry" {
		t.Fatalf("unexpected model: got %q, want %q", resp.Model, "Camry")
	}
}
