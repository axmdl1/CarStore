package usecase

import (
	"CarStore/OrderService/internal/entity"
	_interface "CarStore/OrderService/internal/repository/interface"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type OrderUsecase struct {
	repo _interface.IOrderRepo
	nc   *nats.Conn
}

func NewOrderUsecase(r _interface.IOrderRepo, nc *nats.Conn) *OrderUsecase {
	return &OrderUsecase{repo: r, nc: nc}
}

func (o *OrderUsecase) Create(ctx context.Context, order *entity.Order) error {
	if err := o.repo.Create(ctx, order); err != nil {
		return err
	}

	// publish
	evt := struct {
		OrderID   string    `json:"order_id"`
		CarID     string    `json:"car_id"`
		Quantity  int       `json:"quantity"`
		CreatedAt time.Time `json:"created_at"`
	}{
		OrderID:   order.ID.String(),
		CarID:     order.CarID.String(),
		Quantity:  order.Quantity,
		CreatedAt: order.CreatedAt,
	}
	data, _ := json.Marshal(evt)
	if err := o.nc.Publish("order.created", data); err != nil {
		log.Printf("warning: failed to publish order.created: %v", err)
	} else {
		log.Printf("published order.created for order %s", order.ID)
	}

	return nil
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
