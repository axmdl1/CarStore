package repository

import (
	"CarStore/UserService/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	SetVerificationCode(ctx context.Context, email, code string, expires time.Time) error
	VerifyCode(ctx context.Context, email, code string) (*entity.User, error)
	ChangeRole(ctx context.Context, id, role string) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
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
	update := bson.M{"$set": bson.M{
		"is_active":    user.IsActive,
		"verif_code":   user.VerificationCode,
		"code_expires": user.CodeExpiresAt,
	}}
	_, err := u.collection.UpdateOne(ctx,
		bson.M{"id": user.ID},
		update,
	)
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

func (u *userRepositoryMongo) SetVerificationCode(ctx context.Context, email, code string, expires time.Time) error {
	_, err := u.collection.UpdateOne(ctx,
		bson.M{"email": email},
		bson.M{"$set": bson.M{"verif_code": code, "code_expires": expires}},
	)
	return err
}

func (u *userRepositoryMongo) VerifyCode(ctx context.Context, email, code string) (*entity.User, error) {
	var user entity.User
	err := u.collection.FindOne(ctx, bson.M{
		"email":        email,
		"verif_code":   code,
		"code_expires": bson.M{"$gte": time.Now()},
	}).Decode(&user)
	return &user, err
}

func (u *userRepositoryMongo) ChangeRole(ctx context.Context, id, role string) (*entity.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %v", err)
	}

	filter := bson.M{"id": uid}
	update := bson.M{"$set": bson.M{"role": role}}
	after := options.After
	res := u.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(after))
	var user entity.User
	if err := res.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepositoryMongo) DeleteUser(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user id: %v", err)
	}

	filter := bson.M{"id": uid}
	res, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
