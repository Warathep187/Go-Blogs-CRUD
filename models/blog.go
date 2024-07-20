package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Blog struct {
	ID        string             `bson:"_id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	CreatedBy string             `json:"createdBy"`
	CreatedAt primitive.DateTime `json:"createdAt" swaggertype:"string"`
}
