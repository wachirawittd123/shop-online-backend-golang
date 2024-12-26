package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a user document in the database
type Product struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name" binding:"required"`
	Price      int                `bson:"price" binding:"required"`
	Detail     string             `bson:"detail"`
	IDCategory primitive.ObjectID `bson:"id_category,omitempty"`
	CreatedAt  primitive.DateTime `bson:"created_at"`
	UpdatedOn  primitive.DateTime `bson:"updated_on"`
}
