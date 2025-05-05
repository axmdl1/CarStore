package usecase

import (
	"CarStore/CarService/internal/entity"
	_interface "CarStore/CarService/internal/repository/interface"
	"context"
	"github.com/google/uuid"
)

type CarUsecase struct {
	repo _interface.CarRepo
}

func NewCarUsecase(r _interface.CarRepo) *CarUsecase {
	return &CarUsecase{repo: r}
}

func (uc *CarUsecase) Create(ctx context.Context, car *entity.Car) error {
	return uc.repo.Create(ctx, car)
}

func (uc *CarUsecase) GetByID(ctx context.Context, id string) (*entity.Car, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *CarUsecase) Update(ctx context.Context, car *entity.Car) error {
	return uc.repo.Update(ctx, car)
}

func (uc *CarUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *CarUsecase) List(ctx context.Context) ([]*entity.Car, error) {
	return uc.repo.List(ctx)
}

func (u *CarUsecase) DecreaseStock(ctx context.Context, id string, qty int) (int, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return 0, err
	}
	return u.repo.DecreaseStock(ctx, uid, qty)
}
