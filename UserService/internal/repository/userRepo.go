package repository

import (
	"CarStore/UserService/internal/entity"
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
}

type userRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepositoryMongo{
		collection: db.Collection("users"),
	}
}

func (u userRepositoryMongo) Create(ctx context.Context, user *entity.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	_, err := u.collection.InsertOne(ctx, user)
	return err
}

func (u userRepositoryMongo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepositoryMongo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := u.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u userRepositoryMongo) Update(ctx context.Context, user *entity.User) error {
	filter := bson.M{"id": user.ID}
	_, err := u.collection.UpdateOne(ctx, filter, user)
	return err
}

func (u userRepositoryMongo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	uid, _ := uuid.Parse(id)
	err := u.collection.FindOne(ctx, bson.M{"id": uid}).Decode(&user)
	return &user, err
}

func (u userRepositoryMongo) FindAll(ctx context.Context) ([]*entity.User, error) {
	cursor, err := u.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*entity.User
	for cursor.Next(ctx) {
		var user entity.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
