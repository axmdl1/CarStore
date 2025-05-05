package _interface

import (
	"CarStore/CarService/internal/entity"
	"context"
	"github.com/google/uuid"
)

type CarRepo interface {
	Create(ctx context.Context, car *entity.Car) error
	Update(ctx context.Context, car *entity.Car) error
	GetByID(ctx context.Context, id string) (*entity.Car, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entity.Car, error)
	DecreaseStock(ctx context.Context, id uuid.UUID, qty int) (int, error)
}
