package entity

import (
	"github.com/google/uuid"
	"time"
)

type Car struct {
	ID             uuid.UUID `json:"id" bson:"id"`
	Brand          string    `json:"brand" bson:"brand"`
	Model          string    `json:"model" bson:"model"`
	Year           int       `json:"year" bson:"year"`
	Price          float64   `json:"price" bson:"price"`
	Description    string    `json:"description" bson:"description"`
	EngineCapacity float64   `json:"engine_capacity" bson:"engine_capacity"`
	Mileage        int       `json:"mileage" bson:"mileage"`
	Gearbox        string    `json:"gearbox" bson:"gearbox"`
	EngineType     string    `json:"engine_type" bson:"engine_type"`
	Stock          int       `json:"stock" bson:"stock"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}
