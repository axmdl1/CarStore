package usecase

import (
	"CarStore/OrderService/internal/entity"
	_interface "CarStore/OrderService/internal/repository/interface"
	"context"
)

type OrderUsecase struct {
	repo _interface.IOrderRepo
}

func NewOrderUsecase(r _interface.IOrderRepo) *OrderUsecase {
	return &OrderUsecase{repo: r}
}

func (o *OrderUsecase) Create(ctx context.Context, order *entity.Order) error {
	return o.repo.Create(ctx, order)
}

func (o *OrderUsecase) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	return o.repo.GetByID(ctx, id)
}

func (o *OrderUsecase) Update(ctx context.Context, order *entity.Order) error {
	return o.repo.Update(ctx, order)
}

func (o *OrderUsecase) Delete(ctx context.Context, id string) error {
	return o.repo.Delete(ctx, id)
}

func (o *OrderUsecase) List(ctx context.Context) ([]*entity.Order, error) {
	return o.repo.List(ctx)
}
