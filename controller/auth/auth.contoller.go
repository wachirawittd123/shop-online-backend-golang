package authController

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
	"go.mongodb.org/mongo-driver/bson"
)

// Login authenticates a user and generates a JWT token
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind request data
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	collection, ctx := common.GetCollection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Validate password
	if !common.ComparePasswords(user.Password, loginData.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := common.GenerateToken(user.ID.Hex(), user.Role)
	log.Println(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"token": token}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Logout(c *gin.Context) {
	// Get the token from the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	// Parse the token to get expiration time
	claims, err := common.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Add the token to the blacklist
	expirationTime := time.Unix(claims.ExpiresAt, 0)
	common.AddToBlacklist(tokenString, expirationTime)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
