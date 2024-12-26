package userController

import (
	"context"
	"log"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
	log.Println("filter===========>")
	log.Println(filter)

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

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// AddUser adds a new user to the database
func AddUser(c *gin.Context) {
	var user models.User

	// Bind JSON to the User struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection, ctx := common.GetCollection("users")

	// Set role with validation
	errRole := user.SetRole(user.Role)
	if errRole != "" {
		c.JSON(http.StatusConflict, gin.H{"error": errRole})
		return
	}

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		// Email already exists
		c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
		return
	}

	hashedPassword, err := common.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	// Set additional fields
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Insert the user into the database
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
}

func RemoveUser(c *gin.Context) {
	// Get the user ID from the request
	_id := c.Param("id")

	// Convert the ID to ObjectID
	objectID, err := common.ConvertIDMongodb(_id, c)

	if err != nil {
		return
	}

	collection, ctx := common.GetCollection("users")

	// Perform the deletion
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	// Check if a user was deleted
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUserName updates the name of a user by their ID
func UpdateUser(c *gin.Context) {
	// Get the user ID from the URL parameter
	_id := c.Param("id")

	// Parse the user ID as a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Parse the request body to get the new name
	var requestBody struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection, ctx := common.GetCollection("users")

	// Perform the update
	update := bson.M{"$set": bson.M{"name": requestBody.Name}}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user name"})
		return
	}

	// Check if a user was actually updated
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "User name updated successfully"})
}
