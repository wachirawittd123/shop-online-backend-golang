package cartController

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func pilineQuery(baseMatch bson.D, search string) mongo.Pipeline {
	pipeline := mongo.Pipeline{}

	// Add base match stage if it's not empty
	if len(baseMatch) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: baseMatch}})
	}

	// Unwind items
	pipeline = append(pipeline, bson.D{
		{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$items"}, {Key: "preserveNullAndEmptyArrays", Value: true}}},
	})

	// Lookup products
	pipeline = append(pipeline, bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "products"},
			{Key: "localField", Value: "items.product_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "product_details"},
		}},
	})

	// Unwind product_details
	pipeline = append(pipeline, bson.D{
		{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$product_details"}, {Key: "preserveNullAndEmptyArrays", Value: true}}},
	})

	// Add search filter for product name
	if search != "" {
		pipeline = append(pipeline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "product_details.name", Value: bson.M{
					"$regex":   search,
					"$options": "i",
				}},
			}},
		})
	}

	// Group results
	pipeline = append(pipeline, bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "user_id", Value: bson.D{{Key: "$first", Value: "$user_id"}}},
			{Key: "status", Value: bson.D{{Key: "$first", Value: "$status"}}},
			{Key: "sub_total", Value: bson.D{{Key: "$first", Value: "$sub_total"}}},
			{Key: "total", Value: bson.D{{Key: "$first", Value: "$total"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "updated_on", Value: bson.D{{Key: "$first", Value: "$updated_on"}}},
			{Key: "items", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "product_id", Value: "$items.product_id"},
				{Key: "qty", Value: "$items.qty"},
				{Key: "total", Value: "$items.total"},
				{Key: "product_details", Value: "$product_details"},
			}}}},
		}},
	})

	return pipeline
}

func buildMatchStage(args RequestBuildMatchStage) (bson.D, error) {
	filter := bson.D{}

	// Filter by cart ID if provided
	if args.CardID != "" {
		objectCartID, err := primitive.ObjectIDFromHex(args.CardID)
		if err != nil {
			return nil, fmt.Errorf("invalid cart ID format")
		}
		filter = append(filter, bson.E{Key: "_id", Value: objectCartID})
	}

	// Return an empty filter if no conditions are added
	if len(filter) == 0 {
		return bson.D{}, nil
	}

	return filter, nil
}

// getUserIDFromContext retrieves the user ID from the request context
func getUserIDFromContext(c *gin.Context) (primitive.ObjectID, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return primitive.NilObjectID, fmt.Errorf("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("user ID is not a valid string")
	}

	return common.ConvertIDMongodb(userIDStr, c)
}

// getCartID retrieves the cart ID from the request or the active cart for the user
func getCartID(userID primitive.ObjectID, requestCartID string) (string, error) {
	collection, ctx := common.GetCollection("carts")
	defer ctx.Done()

	var existingCart models.Cart
	err := collection.FindOne(ctx, bson.M{"user_id": userID, "status": "active"}).Decode(&existingCart)

	if requestCartID != "" {
		return requestCartID, nil
	} else if existingCart.ID.Hex() != "000000000000000000000000" {
		return existingCart.ID.Hex(), nil
	}
	return "", err
}

// updateCart updates an existing cart in the database
func updateCart(cartID string, requestBody RequestUpdateCart, cartItems []models.CartItem, c *gin.Context) error {
	objectID, err := common.ConvertIDMongodb(cartID, c)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"total":      requestBody.Total,
			"sub_total":  requestBody.SubTotal,
			"items":      cartItems,
			"updated_on": primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	return common.UpdateOneCommonInDB(objectID, update, c, "carts")
}

// createCart creates a new cart in the database
func createCart(userID primitive.ObjectID, requestBody RequestUpdateCart, cartItems []models.CartItem) error {
	newCart := models.Cart{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Status:    "active",
		Total:     requestBody.Total,
		SubTotal:  requestBody.SubTotal,
		Items:     cartItems,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedOn: primitive.NewDateTimeFromTime(time.Now()),
	}

	return insertCart(newCart)
}

// mapRequestItemsToCartItems maps request items to cart items with valid MongoDB ObjectIDs
func mapRequestItemsToCartItems(items []RequestItemsCart) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	for _, item := range items {
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("invalid product ID: %s", item.ProductID)
		}
		cartItems = append(cartItems, models.CartItem{
			ProductID: productID,
			Qty:       item.Qty,
			Total:     item.Total,
		})
	}
	return cartItems, nil
}

// insertCart inserts a new cart into the database
func insertCart(cart models.Cart) error {
	collection, ctx := common.GetCollection("carts")
	defer ctx.Done()

	_, err := collection.InsertOne(ctx, cart)
	return err
}
