package cache_test

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

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
	return &fakeRepo{
		cars: map[string]models.Car{},
	}
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
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got1 == nil || got1.ID != car.ID {
		t.Fatalf("want id=%s, got=%v", car.ID, got1)
	}
	if repo.calls.get != 1 {
		t.Fatalf("repo get calls=%d, want=1", repo.calls.get)
	}
	got2, err := c.GetCarByID(ctx, "id1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if repo.calls.get != 1 {
		t.Fatalf("repo get calls=%d, want still=1 (cache hit)", repo.calls.get)
	}
	if got2 == nil || got2.ID != got1.ID {
		t.Fatalf("cache returned different car")
	}
}

func TestCarCache_GetCarByID_ExpiredTTL(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	car := models.Car{ID: "id1", Brand: "Audi", Model: "A6", Year: 2021}
	repo.cars[car.ID] = car

	c := cache.NewCarCache(repo, 30*time.Millisecond)
	if _, err := c.GetCarByID(ctx, car.ID); err != nil {
		t.Fatal(err)
	}
	if repo.calls.get != 1 {
		t.Fatalf("want repo get=1, got=%d", repo.calls.get)
	}
	time.Sleep(40 * time.Millisecond)
	if _, err := c.GetCarByID(ctx, car.ID); err != nil {
		t.Fatal(err)
	}
	if repo.calls.get != 2 {
		t.Fatalf("want repo get=2 after ttl, got=%d", repo.calls.get)
	}
}

func TestCarCache_ListCars_CachesOnce(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	expected := []models.Car{
		{ID: "1", Brand: "Toyota", Model: "Camry", Year: 2020},
	}
	repo.list = expected

	c := cache.NewCarCache(repo, time.Minute)
	got1, err := c.ListCars(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, got1) {
		t.Fatalf("want=%v, got=%v", expected, got1)
	}
	if repo.calls.list != 1 {
		t.Fatalf("repo list calls=%d, want=1", repo.calls.list)
	}
	got2, err := c.ListCars(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if repo.calls.list != 1 {
		t.Fatalf("repo list calls=%d, want still=1 (cache hit)", repo.calls.list)
	}
	if len(got2) > 0 {
		got2[0].Brand = "MUTATED"
	}
	got3, err := c.ListCars(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got3[0].Brand != expected[0].Brand {
		t.Fatalf("cache slice must be cloned; want=%s, got=%s", expected[0].Brand, got3[0].Brand)
	}
}

func TestCarCache_Insert_Update_Delete_Invalidates(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	c := cache.NewCarCache(repo, time.Minute)
	repo.list = []models.Car{}
	if _, err := c.ListCars(ctx); err != nil {
		t.Fatal(err)
	}
	if repo.calls.list != 1 {
		t.Fatalf("want repo list=1, got=%d", repo.calls.list)
	}
	newCar := models.Car{ID: "id2", Brand: "Honda", Model: "Civic", Year: 2019}
	if err := c.InsertCar(ctx, &newCar); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListCars(ctx); err != nil {
		t.Fatal(err)
	}
	if repo.calls.list != 2 {
		t.Fatalf("want repo list=2 after insert, got=%d", repo.calls.list)
	}

	newCar.Brand = "HONDA"
	if err := c.UpdateCar(ctx, &newCar); err != nil {
		t.Fatal(err)
	}

	repo.calls.get = 0
	if _, err := c.GetCarByID(ctx, "id2"); err != nil {
		t.Fatal(err)
	}

	if _, err := c.GetCarByID(ctx, "id2"); err != nil {
		t.Fatal(err)
	}
	if repo.calls.get > 1 {
		t.Fatalf("want repo get <=1 after update (warm or first miss), got=%d", repo.calls.get)
	}

	if err := c.DeleteByID(ctx, "id2"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.ListCars(ctx); err != nil {
		t.Fatal(err)
	}
	if repo.calls.list != 3 {
		t.Fatalf("want repo list=3 after delete, got=%d", repo.calls.list)
	}
}

func TestCarCache_ErrorsAreNotCached(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := newFakeRepo()
	repo.err = apperr.ErrNotFound

	c := cache.NewCarCache(repo, time.Minute)
	if _, err := c.GetCarByID(ctx, "nope"); err == nil {
		t.Fatalf("want error, got nil")
	}
	getCalls := repo.calls.get
	if _, err := c.GetCarByID(ctx, "nope"); err == nil {
		t.Fatalf("want error again, got nil")
	}
	if repo.calls.get != getCalls+1 {
		t.Fatalf("want repo get increment on repeat error, got=%d", repo.calls.get)
	}
}
