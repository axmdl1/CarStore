package usecase

import (
	"CarStore/UserService/internal/entity"
	"CarStore/UserService/internal/repository"
	"CarStore/UserService/pkg/email"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

type JWTService interface {
	GenerateToken(userID string, role string) (string, error)
}

type UserUsecase struct {
	repo        repository.UserRepository
	jwtSvc      JWTService
	emailSender email.Sender
	rdb         *redis.Client
}

func NewUserUsecase(r repository.UserRepository, j JWTService, e email.Sender, rdb *redis.Client) *UserUsecase {
	return &UserUsecase{repo: r, jwtSvc: j, emailSender: e, rdb: rdb}
}

func (u *UserUsecase) Register(ctx context.Context, email, username, password, role string) error {
	//Checking for unique username and email
	if _, err := u.repo.FindByEmail(ctx, email); err == nil {
		return errors.New("email already in use")
	}
	if _, err := u.repo.FindByUsername(ctx, username); err == nil {
		return errors.New("username already in use")
	}

	//hashing password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &entity.User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  string(hashed),
		Role:      role,
		CreatedAt: time.Now().UTC(),
		IsActive:  false,
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return err
	}

	_, err = u.SendVerificationCode(ctx, email)
	return err

	/*token, err := u.jwtSvc.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}
	return token, nil*/
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

func (u *UserUsecase) Profile(ctx context.Context, id string) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("user id required")
	}
	key := "user:profile:" + id

	data, err := u.rdb.Get(ctx, key).Bytes()
	if err == nil {
		var cached entity.User
		if err := json.Unmarshal(data, &cached); err == nil {
			return &cached, nil
		}
	}

	user, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// cache result
	if buf, err := json.Marshal(user); err == nil {
		u.rdb.Set(ctx, key, buf, 10*time.Minute)
	}
	return user, nil
}

func (u *UserUsecase) List(ctx context.Context) ([]*entity.User, error) {
	return u.repo.FindAll(ctx)
}

func (u *UserUsecase) SendVerificationCode(ctx context.Context, email string) (string, error) {
	code := fmt.Sprintf("%06d", rand.Intn(1_000_000))
	expires := time.Now().Add(15 * time.Minute)
	if err := u.repo.SetVerificationCode(ctx, email, code, expires); err != nil {
		return "", err
	}
	body := fmt.Sprintf("Your verification code is %s. It expires in 15m.", code)
	if err := u.emailSender.Send(email, "Verify your account", body); err != nil {
		return "", err
	}

	// invalidate cache if exists
	if u.rdb != nil {
		u.rdb.Del(ctx, fmt.Sprintf("user:profile:%s", email))
	}
	return "code_sent", nil
}

func (u *UserUsecase) ConfirmEmail(ctx context.Context, email, code string) (string, error) {
	user, err := u.repo.VerifyCode(ctx, email, code)
	if err != nil {
		return "", err
	}
	user.IsActive = true
	user.VerificationCode = ""
	user.CodeExpiresAt = time.Time{}
	if err := u.repo.Update(ctx, user); err != nil {
		return "", err
	}
	token, err := u.jwtSvc.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		return "", err
	}
	u.rdb.Del(ctx, fmt.Sprintf("user:profile:%s", user.ID.String()))
	return token, nil
}
