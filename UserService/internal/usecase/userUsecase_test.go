package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"CarStore/UserService/internal/entity"
	"CarStore/UserService/pkg/email"
	"CarStore/UserService/pkg/jwt"
)

// mockRepo implements repository.UserRepository for testing.
type mockRepo struct {
	user *entity.User
	err  error
}

func (m *mockRepo) Create(ctx context.Context, u *entity.User) error {
	return nil
}
func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.user, m.err
}
func (m *mockRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.user, m.err
}
func (m *mockRepo) Update(ctx context.Context, u *entity.User) error {
	m.user = u
	return nil
}
func (m *mockRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	return m.user, m.err
}
func (m *mockRepo) FindAll(ctx context.Context) ([]*entity.User, error) {
	return []*entity.User{m.user}, nil
}
func (m *mockRepo) SetVerificationCode(ctx context.Context, email, code string, expires time.Time) error {
	m.user.VerificationCode = code
	m.user.CodeExpiresAt = expires
	return nil
}
func (m *mockRepo) VerifyCode(ctx context.Context, email, code string) (*entity.User, error) {
	if m.user.VerificationCode == code && time.Now().Before(m.user.CodeExpiresAt) {
		return m.user, nil
	}
	return nil, m.err
}

// setupUsecaseWithRedis returns a UserUsecase wired with an in-memory Redis and a stub repo,
// plus the stub user's ID for testing.
func setupUsecaseWithRedis(t *testing.T) (*UserUsecase, *miniredis.Miniredis, string) {
	t.Helper()

	// start in-memory Redis
	mredis, err := miniredis.Run()
	assert.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr: mredis.Addr(),
	})

	jwtSvc := jwt.NewJWTService("testsecret", "UserService")
	emailSvc := email.NewConsoleSender()

	// stub user with a real UUID
	stubID := uuid.New()
	stubUser := &entity.User{
		ID:               stubID,
		Email:            "a@b.com",
		Username:         "u1",
		IsActive:         true,
		VerificationCode: "",
	}
	repo := &mockRepo{user: stubUser, err: nil}

	uc := NewUserUsecase(repo, jwtSvc, emailSvc, rdb)
	return uc, mredis, stubID.String()
}

func TestProfile_CacheAside(t *testing.T) {
	uc, mredis, stubID := setupUsecaseWithRedis(t)
	defer mredis.Close()

	ctx := context.Background()
	key := "user:profile:" + stubID

	// 1) First call: cache miss, fetch from repo, then cache
	user1, err := uc.Profile(ctx, stubID)
	assert.NoError(t, err)
	assert.Equal(t, "a@b.com", user1.Email)
	assert.True(t, mredis.Exists(key))

	// 2) Change repo behind the scenes
	uc.repo.(*mockRepo).user.Email = "changed@b.com"

	// 2nd call: still returns cached old value
	user2, err := uc.Profile(ctx, stubID)
	assert.NoError(t, err)
	assert.Equal(t, "a@b.com", user2.Email)

	// 3) Evict cache and retry: returns updated value
	mredis.Del(key)
	user3, err := uc.Profile(ctx, stubID)
	assert.NoError(t, err)
	assert.Equal(t, "changed@b.com", user3.Email)
}

func TestRegister_Validation(t *testing.T) {
	uc, mredis, _ := setupUsecaseWithRedis(t)
	defer mredis.Close()

	ctx := context.Background()

	// missing email
	err := uc.Register(ctx, "", "u1", "p1", "user")
	assert.Error(t, err)

	// missing username
	err = uc.Register(ctx, "a@b.com", "", "p1", "user")
	assert.Error(t, err)

	// missing password
	err = uc.Register(ctx, "a@b.com", "u1", "", "user")
	assert.Error(t, err)
}
