package deliveryController

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestCreateDelivery struct {
	OrderID     primitive.ObjectID `json:"order_id"`     // Reference to the cart/order ID
	UserID      primitive.ObjectID `json:"user_id"`      // Reference to the user
	DeliveryFee float64            `json:"delivery_fee"` // Cost of delivery
}
