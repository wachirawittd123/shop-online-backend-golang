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

var user struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      string             `bson:"role"`
	Token     string             `bson:"token"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

func GetUsers(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	collection := common.GetUsersCollection()
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

	// Check if email already exists
	collection := common.GetUsersCollection()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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