package repository

import (
	"CarStore/CarService/internal/entity"
	_interface "CarStore/CarService/internal/repository/interface"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type carRepo struct {
	coll *mongo.Collection
}

func NewCarRepo(db *mongo.Database) _interface.CarRepo {
	return &carRepo{coll: db.Collection("cars")}
}

func (c carRepo) Create(ctx context.Context, car *entity.Car) error {
	if car.ID == uuid.Nil {
		car.ID = uuid.New()
	}
	car.CreatedAt = time.Now()
	_, err := c.coll.InsertOne(ctx, car)
	return err
}

func (c carRepo) Update(ctx context.Context, car *entity.Car) error {
	_, err := c.coll.ReplaceOne(ctx, bson.M{"id": car.ID}, car)
	return err
}

func (c carRepo) GetByID(ctx context.Context, id string) (*entity.Car, error) {
	var car entity.Car
	uid, _ := uuid.Parse(id)
	err := c.coll.FindOne(ctx, bson.M{"id": uid}).Decode(&car)
	return &car, err
}

func (c carRepo) Delete(ctx context.Context, id string) error {
	uid, _ := uuid.Parse(id)
	_, err := c.coll.DeleteOne(ctx, bson.M{"id": uid})
	return err
}

func (c carRepo) List(ctx context.Context) ([]*entity.Car, error) {
	cursor, err := c.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []*entity.Car
	for cursor.Next(ctx) {
		var c entity.Car
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		cars = append(cars, &c)
	}

	return cars, nil
}

func (c carRepo) DecreaseStock(ctx context.Context, id uuid.UUID, qty int) (int, error) {
	res := c.coll.FindOneAndUpdate(ctx,
		bson.M{"id": id, "stock": bson.M{"$gte": qty}},
		bson.M{"$inc": bson.M{"stock": -qty}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var updated entity.Car
	if err := res.Decode(&updated); err != nil {
		return 0, err
	}
	return updated.Stock, nil
}
