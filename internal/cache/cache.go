package cache

import (
	"context"
	"sync"
	"time"

	"github.com/pavel97go/service-cars/internal/apperr"
	"github.com/pavel97go/service-cars/internal/models"
	"github.com/pavel97go/service-cars/internal/repository"
)

type byIDItem struct {
	car models.Car
	exp time.Time
}
type listItem struct {
	cars []models.Car
	exp  time.Time
}
type CarCache struct {
	next repository.CarProvider
	mu   sync.RWMutex
	ttl  time.Duration
	id   map[string]byIDItem
	list *listItem
}

func NewCarCache(next repository.CarProvider, ttl time.Duration) *CarCache {
	return &CarCache{
		next: next,
		ttl:  ttl,
		id:   make(map[string]byIDItem),
	}
}
func (c *CarCache) getNow() time.Time {
	return time.Now()
}
func cloneCars(in []models.Car) []models.Car {
	if len(in) == 0 {
		return nil
	}
	out := make([]models.Car, len(in))
	copy(out, in)
	return out
}
func (c *CarCache) getByID(id string) (models.Car, bool) {
	c.mu.RLock()
	item, ok := c.id[id]
	if !ok {
		c.mu.RUnlock()
		return models.Car{}, false
	}
	if c.getNow().After(item.exp) {
		c.mu.RUnlock()
		return models.Car{}, false
	}
	car := item.car
	c.mu.RUnlock()
	return car, true
}
func (c *CarCache) setByID(car models.Car) {
	c.mu.Lock()
	c.id[car.ID] = byIDItem{
		car: car,
		exp: c.getNow().Add(c.ttl),
	}
	c.mu.Unlock()
}
func (c *CarCache) delByID(id string) {
	c.mu.Lock()
	delete(c.id, id)
	c.mu.Unlock()
}
func (c *CarCache) getList() ([]models.Car, bool) {
	c.mu.RLock()
	item := c.list
	if item == nil {
		c.mu.RUnlock()
		return nil, false
	}
	if c.getNow().After(item.exp) {
		c.mu.RUnlock()
		return nil, false
	}
	out := cloneCars(item.cars)
	c.mu.RUnlock()
	return out, true

}
func (c *CarCache) setList(cars []models.Car) {
	c.mu.Lock()
	c.list = &listItem{
		cars: cloneCars(cars),
		exp:  c.getNow().Add(c.ttl),
	}
	c.mu.Unlock()
}
func (c *CarCache) invalidateList() {
	c.mu.Lock()
	c.list = nil
	c.mu.Unlock()
}
func (c *CarCache) ListCars(ctx context.Context) ([]models.Car, error) {
	if cars, ok := c.getList(); ok {
		return cars, nil
	}
	cars, err := c.next.ListCars(ctx)
	if err != nil {
		return nil, err
	}
	c.setList(cars)
	return cars, nil
}
func (c *CarCache) GetCarByID(ctx context.Context, id string) (*models.Car, error) {
	if car, ok := c.getByID(id); ok {
		cc := car
		return &cc, nil
	}
	car, err := c.next.GetCarByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c.setByID(*car)
	return car, nil
}
func (c *CarCache) InsertCar(ctx context.Context, newCar *models.Car) error {
	if err := c.next.InsertCar(ctx, newCar); err != nil {
		return err
	}
	c.invalidateList()
	c.setByID(*newCar)
	return nil
}
func (c *CarCache) UpdateCar(ctx context.Context, updatedCar *models.Car) error {
	if err := c.next.UpdateCar(ctx, updatedCar); err != nil {
		if err == apperr.ErrNotFound {
			return apperr.ErrNotFound
		}
		return err
	}
	c.delByID(updatedCar.ID)
	c.invalidateList()
	c.setByID(*updatedCar)
	return nil
}
func (c *CarCache) DeleteByID(ctx context.Context, id string) error {
	if err := c.next.DeleteByID(ctx, id); err != nil {
		if err == apperr.ErrNotFound {
			return apperr.ErrNotFound
		}
		return err
	}
	c.delByID(id)
	c.invalidateList()
	return nil
}
