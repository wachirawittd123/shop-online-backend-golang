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

// validateUserRole validates and sets the user's role
func validateUserRole(user *models.User, c *gin.Context) error {
	errRole := user.SetRole(user.Role)
	if errRole != "" {
		common.RespondWithError(c, http.StatusConflict, errRole, nil)
		return nil
	}
	return nil
}

// isEmailInUse checks if the email is already in use
func isEmailInUse(email string, c *gin.Context) bool {
	collection, ctx := common.GetCollection("users")
	defer ctx.Done()

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	if err == nil {
		common.RespondWithError(c, http.StatusConflict, "Email already in use", nil)
		return true
	}
	return false
}

// insertUser sets additional fields and inserts the user into the database
func insertUser(user models.User, c *gin.Context) error {
	collection, ctx := common.GetCollection("users")
	defer ctx.Done()

	// Set additional fields
	user.ID = primitive.NewObjectID()
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to add user", err)
		return err
	}
	return nil
}
