package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebviewServer struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Status    string             `json:"status,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty"`
}
