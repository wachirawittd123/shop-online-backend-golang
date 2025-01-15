package userController

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUsers(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	collection, ctx := common.GetCollection("users")

	filter := bson.M{}
	if search != "" {
		// Filter by search if provided
		filter = bson.M{
			"$or": []bson.M{
				{"email": bson.M{"$regex": search, "$options": "i"}},
				{"name": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	// Parse startDate and endDate if provided
	dateFilter := bson.M{}
	if startDate != "" {
		start, err := time.Parse("2006-01-02", startDate) // Format: YYYY-MM-DD
		if err == nil {
			dateFilter["$gte"] = primitive.NewDateTimeFromTime(start)
		}
	}
	if endDate != "" {
		end, err := time.Parse("2006-01-02", endDate) // Format: YYYY-MM-DD
		if err == nil {
			dateFilter["$lte"] = primitive.NewDateTimeFromTime(end)
		}
	}
	if len(dateFilter) > 0 {
		filter["created_at"] = dateFilter
	}

	// Find all users
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	// Decode users into a slice
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	if users == nil {
		users = []models.User{}
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "status_code": http.StatusOK})
}

// AddUser adds a new user to the database
func AddUser(c *gin.Context) {
	var user models.User

	// Parse and validate the request body
	if err := c.ShouldBindJSON(&user); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate and set the role
	if err := validateUserRole(&user, c); err != nil {
		return
	}

	// Check if the email already exists
	if isEmailInUse(user.Email, c) {
		return
	}

	hashedPassword, err := common.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = hashedPassword

	// Set additional fields and insert the user
	if err := insertUser(user, c); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "status_code": http.StatusOK})
}

func RemoveUser(c *gin.Context) {
	objectID, err := common.ConvertIDMongodb(c.Param("id"), c)
	if err != nil {
		return
	}

	if err := common.DeleteOneCommonByID(objectID, c, "users"); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "status_code": http.StatusOK})
}

// UpdateUserName updates the name of a user by their ID
func UpdateUser(c *gin.Context) {
	objectID, err := common.ConvertIDMongodb(c.Param("id"), c)
	if err != nil {
		return
	}

	var requestBody RequestUpdateBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	update := bson.M{"$set": bson.M{
		"name": requestBody.Name,
		"shipping_address": bson.M{
			"phone":       requestBody.ShippingAddr.Phone,
			"street":      requestBody.ShippingAddr.Street,
			"city":        requestBody.ShippingAddr.City,
			"state":       requestBody.ShippingAddr.State,
			"postal_code": requestBody.ShippingAddr.PostalCode,
			"country":     requestBody.ShippingAddr.Country,
			"latitude":    requestBody.ShippingAddr.Latitude, // Example latitude
			"longitude":   requestBody.ShippingAddr.Longitude,
		},
	}}

	if err := common.UpdateOneCommonInDB(objectID, update, c, "users"); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User successfully", "status_code": http.StatusOK})
}
