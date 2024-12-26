package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductCategory represents a user document in the database
type ProductCategory struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedOn primitive.DateTime `bson:"updated_on"`
}
