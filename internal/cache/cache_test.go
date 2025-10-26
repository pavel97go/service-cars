package cache_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pavel97go/service-cars/internal/apperr"
	"github.com/pavel97go/service-cars/internal/cache"
	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/repository"
)

type fakeRepo struct {
	mu    sync.Mutex
	cars  map[string]models.Car
	list  []models.Car
	err   error
	calls struct {
		list, get, ins, upd, del int
	}
}

var _ repository.CarProvider = (*fakeRepo)(nil)

func newFakeRepo() *fakeRepo {
	return &fakeRepo{cars: map[string]models.Car{}}
}

func (f *fakeRepo) ListCars(ctx context.Context) ([]models.Car, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls.list++
	if f.err != nil {
		return nil, f.err
	}
	out := make([]models.Car, len(f.list))
	copy(out, f.list)
	return out, nil
}

func (f *fakeRepo) GetCarByID(ctx context.Context, id string) (*models.Car, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls.get++
	if f.err != nil {
		return nil, f.err
	}
	c, ok := f.cars[id]
	if !ok {
		return nil, apperr.ErrNotFound
	}
	cc := c
	return &cc, nil
}

func (f *fakeRepo) InsertCar(ctx context.Context, c *models.Car) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls.ins++
	if f.err != nil {
		return f.err
	}
	f.cars[c.ID] = *c
	f.list = append(f.list, *c)
	return nil
}

func (f *fakeRepo) UpdateCar(ctx context.Context, c *models.Car) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls.upd++
	if f.err != nil {
		return f.err
	}
	if _, ok := f.cars[c.ID]; !ok {
		return apperr.ErrNotFound
	}
	f.cars[c.ID] = *c
	for i := range f.list {
		if f.list[i].ID == c.ID {
			f.list[i] = *c
			break
		}
	}
	return nil
}

func (f *fakeRepo) DeleteByID(ctx context.Context, id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls.del++
	if f.err != nil {
		return f.err
	}
	if _, ok := f.cars[id]; !ok {
		return apperr.ErrNotFound
	}
	delete(f.cars, id)
	out := make([]models.Car, 0, len(f.list))
	for _, c := range f.list {
		if c.ID != id {
			out = append(out, c)
		}
	}
	f.list = out
	return nil
}

func TestCarCache_GetCarByID_MissThenHit(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	car := models.Car{ID: "id1", Brand: "BMW", Model: "X5", Year: 2022}
	repo.cars[car.ID] = car

	c := cache.NewCarCache(repo, 200*time.Millisecond)
	got1, err := c.GetCarByID(ctx, "id1")
	require.NoError(t, err)
	require.NotNil(t, got1)
	assert.Equal(t, car.ID, got1.ID)
	assert.Equal(t, 1, repo.calls.get, "first call should hit repo")

	got2, err := c.GetCarByID(ctx, "id1")
	require.NoError(t, err)
	assert.Equal(t, 1, repo.calls.get, "second call should be cache hit")
	assert.Equal(t, got1.ID, got2.ID)
}

func TestCarCache_GetCarByID_ExpiredTTL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	car := models.Car{ID: "id1", Brand: "Audi", Model: "A6", Year: 2021}
	repo.cars[car.ID] = car

	c := cache.NewCarCache(repo, 30*time.Millisecond)
	_, err := c.GetCarByID(ctx, car.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, repo.calls.get)

	time.Sleep(50 * time.Millisecond)
	_, err = c.GetCarByID(ctx, car.ID)
	require.NoError(t, err)
	assert.Equal(t, 2, repo.calls.get, "should call repo again after TTL")
}

func TestCarCache_ListCars_CachesOnce(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	expected := []models.Car{{ID: "1", Brand: "Toyota", Model: "Camry", Year: 2020}}
	repo.list = expected

	c := cache.NewCarCache(repo, time.Minute)
	got1, err := c.ListCars(ctx)
	require.NoError(t, err)
	assert.Equal(t, expected, got1)
	assert.Equal(t, 1, repo.calls.list)

	got2, err := c.ListCars(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, repo.calls.list, "second list should be cache hit")

	// Проверяем, что кэш возвращает копию
	got2[0].Brand = "MUTATED"
	got3, err := c.ListCars(ctx)
	require.NoError(t, err)
	assert.Equal(t, expected[0].Brand, got3[0].Brand, "cache should not be mutated")
}

func TestCarCache_Insert_Update_Delete_Invalidates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	c := cache.NewCarCache(repo, time.Minute)
	repo.list = []models.Car{}

	_, err := c.ListCars(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, repo.calls.list)

	newCar := models.Car{ID: "id2", Brand: "Honda", Model: "Civic", Year: 2019}
	require.NoError(t, c.InsertCar(ctx, &newCar))

	_, err = c.ListCars(ctx)
	require.NoError(t, err)
	assert.Equal(t, 2, repo.calls.list, "list invalidated after insert")

	newCar.Brand = "HONDA"
	require.NoError(t, c.UpdateCar(ctx, &newCar))

	repo.calls.get = 0
	_, err = c.GetCarByID(ctx, "id2")
	require.NoError(t, err)
	_, err = c.GetCarByID(ctx, "id2")
	require.NoError(t, err)
	assert.LessOrEqual(t, repo.calls.get, 1, "should be cached after first fetch")

	require.NoError(t, c.DeleteByID(ctx, "id2"))
	_, err = c.ListCars(ctx)
	require.NoError(t, err)
	assert.Equal(t, 3, repo.calls.list, "list invalidated after delete")
}

func TestCarCache_ErrorsAreNotCached(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	repo.err = apperr.ErrNotFound

	c := cache.NewCarCache(repo, time.Minute)
	_, err := c.GetCarByID(ctx, "nope")
	require.Error(t, err)

	getCalls := repo.calls.get
	_, err = c.GetCarByID(ctx, "nope")
	require.Error(t, err)
	assert.Equal(t, getCalls+1, repo.calls.get, "errors must not be cached")
}
