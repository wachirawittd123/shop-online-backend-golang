package cartController

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	deliveryController "github.com/wachirawittd123/shop-online-backend-golang/controller/delivery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCart(c *gin.Context) {
	// Parse and validate the cart ID
	cardId := c.Param("id")
	search := c.Query("search")

	var requestParams RequestBuildMatchStage

	if cardId != "" {
		requestParams.CardID = cardId
	}
	matchStage, err := buildMatchStage(requestParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection, ctx := common.GetCollection("carts")
	defer ctx.Done()

	// Aggregation pipeline
	pipeline := pilineQuery(matchStage, search)

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart", "details": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Decode the result into a Cart with products
	var cartWithProducts []bson.M
	if err := cursor.All(ctx, &cartWithProducts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cart data", "details": err.Error()})
		return
	}

	if len(cardId) > 0 {
		c.JSON(http.StatusOK, gin.H{"cart": cartWithProducts[0], "status_code": http.StatusOK})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cartWithProducts, "status_code": http.StatusOK})
}

func UpdateCart(c *gin.Context) {
	// Retrieve and validate the user ID
	userID, err := getUserIDFromContext(c)
	if err != nil {
		common.RespondWithError(c, http.StatusUnauthorized, "User not authenticated", err)
		return
	}

	// Parse and validate the request body
	var requestBody RequestUpdateCart
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Determine the cart ID to use (new or existing)
	newIdCart := getCartID(userID, requestBody.ID)

	// Map requestBody.Items to cart items with valid MongoDB ObjectIDs
	cartItems, err := mapRequestItemsToCartItems(requestBody.Items)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid product ID in items", err)
		return
	}

	var newObjectID primitive.ObjectID
	// Handle update or create logic
	if newIdCart != "" {
		err, id := updateCart(newIdCart, requestBody, cartItems, c)
		if err != nil {
			common.RespondWithError(c, http.StatusInternalServerError, "Failed to update cart", err)
			return
		}
		newObjectID = id
	} else {
		requestBody.Status = "active"
		err, id := createCart(userID, requestBody, cartItems)
		if err != nil {
			common.RespondWithError(c, http.StatusInternalServerError, "Failed to create cart", err)
			return
		}
		newObjectID = id
	}

	if requestBody.Status == "completed" {
		deliveryRequest := deliveryController.RequestCreateDelivery{
			OrderID:     newObjectID,             // Ensure correct type
			UserID:      userID,                  // User ID
			DeliveryFee: requestBody.DeliveryFee, // Delivery fee from request body
		}
		status := deliveryController.AddDelivery(deliveryRequest)
		if status == 400 {
			common.RespondWithError(c, http.StatusInternalServerError, "Failed to add delivery", err)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully", "status_code": http.StatusOK})
}
