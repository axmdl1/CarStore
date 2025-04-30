package repository

import (
	"CarStore/OrderService/internal/entity"
	_interface "CarStore/OrderService/internal/repository/interface"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type orderRepo struct {
	coll *mongo.Collection
}

func NewOrderRepo(db *mongo.Database) _interface.IOrderRepo {
	return &orderRepo{coll: db.Collection("orders")}
}

func (o orderRepo) Create(ctx context.Context, order *entity.Order) error {
	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	order.CreatedAt = time.Now()
	_, err := o.coll.InsertOne(ctx, order)
	return err
}

func (o orderRepo) Update(ctx context.Context, order *entity.Order) error {
	_, err := o.coll.ReplaceOne(ctx, bson.M{"id": order.ID}, order)
	return err
}

func (o orderRepo) GetByID(ctx context.Context, id string) (*entity.Order, error) {
	var order entity.Order
	uid, _ := uuid.Parse(id)
	err := o.coll.FindOne(ctx, bson.M{"id": uid}).Decode(&order)
	return &order, err
}

func (o orderRepo) Delete(ctx context.Context, id string) error {
	uid, _ := uuid.Parse(id)
	_, err := o.coll.DeleteOne(ctx, bson.M{"id": uid})
	return err
}

func (o orderRepo) List(ctx context.Context) ([]*entity.Order, error) {
	cursor, err := o.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*entity.Order
	for cursor.Next(ctx) {
		var o entity.Order
		if err := cursor.Decode(&o); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, nil
}
