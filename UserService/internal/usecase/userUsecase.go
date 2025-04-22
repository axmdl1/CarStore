package usecase

import (
	"CarStore/UserService/internal/entity"
	"CarStore/UserService/internal/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type JWTService interface {
	GenerateToken(userID string, role string) (string, error)
}

type UserUsecase struct {
	repo   repository.UserRepository
	jwtSvc JWTService
}

func NewUserUsecase(r repository.UserRepository, j JWTService) *UserUsecase {
	return &UserUsecase{repo: r, jwtSvc: j}
}

func (u *UserUsecase) Register(ctx context.Context, email, username, password, role string) (string, error) {
	//Checking for unique username and email
	if _, err := u.repo.FindByEmail(ctx, email); err == nil {
		return "", errors.New("email already in use")
	}
	if _, err := u.repo.FindByUsername(ctx, username); err == nil {
		return "", errors.New("username already in use")
	}

	//hashing password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &entity.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  string(hashed),
		Role:      role,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return "", err
	}

	token, err := u.jwtSvc.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *UserUsecase) Login(ctx context.Context, identifier, password string) (string, error) {
	var user *entity.User
	var err error

	if strings.Contains(identifier, "@") {
		user, err = u.repo.FindByEmail(ctx, identifier)
	} else {
		user, err = u.repo.FindByUsername(ctx, identifier)
	}
	if err != nil {
		return "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return u.jwtSvc.GenerateToken(user.ID.String(), user.Role)
}
