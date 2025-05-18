package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"CarStore/CarService/internal/entity"
)

// memoryCarRepo is an in-memory implementation of CarRepo for integration testing.
type memoryCarRepo struct {
	store map[uuid.UUID]*entity.Car
}

func newMemoryCarRepo() *memoryCarRepo {
	return &memoryCarRepo{store: make(map[uuid.UUID]*entity.Car)}
}

func (m *memoryCarRepo) Create(ctx context.Context, car *entity.Car) error {
	m.store[car.ID] = car
	return nil
}

func (m *memoryCarRepo) GetByID(ctx context.Context, id string) (*entity.Car, error) {
	uID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	car, ok := m.store[uID]
	if !ok {
		return nil, fmt.Errorf("car not found")
	}
	return car, nil
}

func (m *memoryCarRepo) List(ctx context.Context) ([]*entity.Car, error) {
	list := make([]*entity.Car, 0, len(m.store))
	for _, c := range m.store {
		list = append(list, c)
	}
	return list, nil
}

func (m *memoryCarRepo) Update(ctx context.Context, car *entity.Car) error {
	if _, ok := m.store[car.ID]; !ok {
		return fmt.Errorf("car not found")
	}
	m.store[car.ID] = car
	return nil
}

func (m *memoryCarRepo) Delete(ctx context.Context, id string) error {
	uID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	if _, ok := m.store[uID]; !ok {
		return fmt.Errorf("car not found")
	}
	delete(m.store, uID)
	return nil
}

func (m *memoryCarRepo) DecreaseStock(ctx context.Context, id uuid.UUID, qty int) (int, error) {
	car, ok := m.store[id]
	if !ok {
		return 0, fmt.Errorf("car not found")
	}
	if car.Stock < qty {
		return car.Stock, fmt.Errorf("insufficient stock")
	}
	car.Stock -= qty
	m.store[id] = car
	return car.Stock, nil
}

func TestCarUsecase_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryCarRepo()
	uc := NewCarUsecase(repo)

	// Create
	c1 := &entity.Car{ID: uuid.New(), Brand: "A", Model: "X", Stock: 10}
	err := uc.Create(ctx, c1)
	assert.NoError(t, err)

	// GetByID
	got, err := uc.GetByID(ctx, c1.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, c1.Brand, got.Brand)

	// List
	list, err := uc.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	c1.Price = 99.99
	err = uc.Update(ctx, c1)
	assert.NoError(t, err)
	updated, _ := uc.GetByID(ctx, c1.ID.String())
	assert.Equal(t, 99.99, updated.Price)

	// Delete
	err = uc.Delete(ctx, c1.ID.String())
	assert.NoError(t, err)
	_, err = uc.GetByID(ctx, c1.ID.String())
	assert.Error(t, err)
}

func TestCarUsecase_DecreaseStock(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryCarRepo()
	uc := NewCarUsecase(repo)

	c := &entity.Car{ID: uuid.New(), Brand: "B", Model: "Y", Stock: 5}
	repo.Create(ctx, c)

	// Successful decrease
	newStock, err := uc.DecreaseStock(ctx, c.ID.String(), 3)
	assert.NoError(t, err)
	assert.Equal(t, 2, newStock)

	// Insufficient stock
	_, err = uc.DecreaseStock(ctx, c.ID.String(), 10)
	assert.Error(t, err)
}
