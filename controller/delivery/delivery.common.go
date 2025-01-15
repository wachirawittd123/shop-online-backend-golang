package deliveryController

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/wachirawittd123/shop-online-backend-golang/common"
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// generateRandomCode generates a random alphanumeric string of the specified length.
func generateRandomCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}

func insertDelivery(delivery RequestCreateDelivery) error {
	// Get MongoDB collection
	collection, ctx := common.GetCollection("delivery")
	defer ctx.Done()

	filter := bson.M{"order_id": delivery.OrderID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to check OrderID existence: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("duplicate OrderID: a delivery with this OrderID already exists")
	}

	// Generate random tracking code
	trackingCode := generateRandomCode(10)

	// Prepare the Delivery document
	newDelivery := bson.M{
		"_id":           primitive.NewObjectID(),
		"order_id":      delivery.OrderID,
		"user_id":       delivery.UserID,
		"status":        models.StatusPending, // Initial status
		"delivery_fee":  delivery.DeliveryFee,
		"created_at":    primitive.NewDateTimeFromTime(time.Now()),
		"updated_on":    primitive.NewDateTimeFromTime(time.Now()),
		"tracking_code": trackingCode,
	}

	// Insert into MongoDB
	_, err = collection.InsertOne(ctx, newDelivery)
	if err != nil {
		return err
	}
	return nil
}
