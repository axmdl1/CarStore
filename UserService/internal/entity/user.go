package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID `json:"id" bson:"id"`
	Email            string    `json:"email" bson:"email"`
	Username         string    `json:"username" bson:"username"`
	Password         string    `json:"-" bson:"password"`
	Role             string    `json:"role" bson:"role"`
	IsActive         bool      `json:"is_active" bson:"is_active"`
	VerificationCode string    `json:"-" bson:"verif_code"`
	CodeExpiresAt    time.Time `json:"-" bson:"code_expires"`
	CreatedAt        time.Time `json:"created_at" bson:"createdat"`
}
