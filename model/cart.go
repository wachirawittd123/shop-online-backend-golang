package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cart represents a shopping cart document in the database
type Cart struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"user_id,omitempty"`
	Items    []CartItem         `bson:"items"`     // List of items in the cart
	Status   string             `bson:"status"`    // e.g., "active", "completed", "cancelled"
	SubTotal float64            `bson:"sub_total"` // Total of item prices before discounts or taxes
	Total    float64            `bson:"total"`     // Final total after discounts and taxes
	// Discount  float64            `bson:"discount"`   // Total discount applied to the cart
	// Tax       float64            `bson:"tax"`        // Total tax applied to the cart
	CreatedAt primitive.DateTime `bson:"created_at"` // Timestamp when the cart was created
	UpdatedOn primitive.DateTime `bson:"updated_on"` // Timestamp when the cart was last updated
}

// CartItem represents an individual item in the shopping cart
type CartItem struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Qty       int                `bson:"qty"`
	Total     float64            `bson:"total"`
}
