package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pavel97go/service-cars/internal/apperr"
	"github.com/pavel97go/service-cars/internal/models"
)

type CarRepo struct {
	pool *pgxpool.Pool
}

func NewCarRepo(pool *pgxpool.Pool) *CarRepo {
	return &CarRepo{pool: pool}

}
func (r *CarRepo) ListCars(ctx context.Context) ([]models.Car, error) {
	query := `
		SELECT id, brand, model, year, created_at
		FROM cars
		ORDER BY created_at DESC;
		`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cars := []models.Car{}

	for rows.Next() {
		var c models.Car
		err := rows.Scan(&c.ID, &c.Brand, &c.Model, &c.Year, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		cars = append(cars, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cars, nil

}
func (r *CarRepo) GetCarByID(ctx context.Context, id string) (*models.Car, error) {
	const query = `
		SELECT id, brand, model, year, created_at
		FROM cars
		WHERE id = $1;
	`
	row := r.pool.QueryRow(ctx, query, id)

	var c models.Car
	err := row.Scan(&c.ID, &c.Brand, &c.Model, &c.Year, &c.CreatedAt)
	if err == pgx.ErrNoRows {
		return nil, apperr.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CarRepo) InsertCar(ctx context.Context, newCar *models.Car) error {
	const query = `
		INSERT INTO cars (brand, model, year)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`
	return r.pool.QueryRow(ctx, query, newCar.Brand, newCar.Model, newCar.Year).
		Scan(&newCar.ID, &newCar.CreatedAt)
}

func (r *CarRepo) UpdateCar(ctx context.Context, c *models.Car) error {
	query := `
	UPDATE cars
	SET brand=$2,model=$3,year=$4
	WHERE id=$1;
	`
	ct, err := r.pool.Exec(ctx, query, c.ID, c.Brand, c.Model, c.Year)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	slog.Debug("car updated", "id", c.ID, "rows", ct.RowsAffected())
	return nil
}
func (r *CarRepo) DeleteByID(ctx context.Context, id string) error {
	query := `
	DELETE FROM cars WHERE id = $1;
	`
	ct, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	slog.Debug("car deleted", "id", id, "rows", ct.RowsAffected())
	return nil
}
