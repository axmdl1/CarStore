package entity

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID         uuid.UUID `json:"id" bson:"id"`
	UserID     uuid.UUID `json:"userId" bson:"userId"`
	CarID      uuid.UUID `json:"carId" bson:"carId"`
	Quantity   int       `json:"quantity" bson:"quantity"`
	TotalPrice float64   `json:"price" bson:"price"`
	Status     string    `json:"status" bson:"status"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
}
