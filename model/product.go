package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a user document in the database
type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Price     int                `bson:"price"`
	Detail    string             `bson:"detail"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedOn primitive.DateTime `bson:"updated_on"`
}
