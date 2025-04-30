package _interface

import (
	"CarStore/OrderService/internal/entity"
	"context"
)

type IOrderRepo interface {
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id string) (*entity.Order, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*entity.Order, error)
}
