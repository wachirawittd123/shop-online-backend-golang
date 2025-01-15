package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Delivery represents a delivery document in the database
type Delivery struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`           // Unique identifier for the delivery
	OrderID          primitive.ObjectID `bson:"order_id"`                // Reference to the cart/order ID
	UserID           primitive.ObjectID `bson:"user_id"`                 // Reference to the user
	DeliveryPersonID primitive.ObjectID `bson:"delivery_person_id"`      // Reference to the delivery person
	Status           string             `bson:"status"`                  // Delivery status (e.g., "pending", "shipped", "delivered")
	TrackingCode     string             `bson:"tracking_code,omitempty"` // Optional tracking code from the courier service
	DeliveryFee      float64            `bson:"delivery_fee"`            // Cost of delivery
	CreatedAt        primitive.DateTime `bson:"created_at"`              // Timestamp for when the delivery was created
	UpdatedOn        primitive.DateTime `bson:"updated_on"`              // Timestamp for when the delivery was last updated
	// ExpectedDate     primitive.DateTime `bson:"expected_date,omitempty"`  // Expected delivery date
	DeliveredDate primitive.DateTime `bson:"delivered_date"` // Actual delivery date
}

// Predefined delivery status constants
const (
	StatusPending    = "pending"     // Delivery has been created but not yet started
	StatusInProgress = "in_progress" // Delivery is currently underway
	StatusShipped    = "shipped"     // Delivery has been shipped
	StatusDelivered  = "delivered"   // Delivery has been completed successfully
	StatusCancelled  = "cancelled"   // Delivery has been cancelled
	StatusFailed     = "failed"      // Delivery attempt was unsuccessful
)
