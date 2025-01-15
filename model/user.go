package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	return role == RoleUser || role == RoleAdmin
}

// User represents a user document in the database
type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	Role         string             `bson:"role"`
	Token        string             `bson:"token"`
	CreatedAt    primitive.DateTime `bson:"created_at"`
	ShippingAddr ShippingAddress    `bson:"shipping_address"`
}

func (u *User) SetRole(role string) string {
	if !IsValidRole(role) {
		return "invalid role"
	}
	u.Role = role
	return ""
}

type ShippingAddress struct {
	Phone      string  `bson:"phone"`       // Contact phone number
	Street     string  `bson:"street"`      // Street address
	City       string  `bson:"city"`        // City
	State      string  `bson:"state"`       // State/Province
	PostalCode string  `bson:"postal_code"` // Postal code
	Country    string  `bson:"country"`     // Country
	Latitude   float64 `bson:"latitude"`    // GPS latitude
	Longitude  float64 `bson:"longitude"`   // GPS longitude
}
