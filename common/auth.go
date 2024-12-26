package common

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain-text password with a salt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords compares a hashed password with a plain-text password
func ComparePasswords(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

// Claims represents the structure of JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token for the given user ID
func GenerateToken(userID string, role string) (string, error) {
	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour) // 24-hour token validity

	// Create the claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Create the token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var jwtSecret = []byte(AppConfig.JWTSecret)
	// Sign the token with the secret
	return token.SignedString(jwtSecret)
}

// ValidateToken validates the JWT token and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	var jwtSecret = []byte(AppConfig.JWTSecret)
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// Validate the token
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// AuthMiddleware validates the JWT token in requests
func AuthMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Check if the token is blacklisted
		if IsBlacklisted(tokenString) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		objectID, err := ConvertIDMongodb(claims.UserID, c)

		if err != nil {
			// The error and response are already handled in ConvertIDMongodb
			return
		}

		collection, ctx := GetCollection("users")

		var user struct {
			Token string `bson:"token"`
		}

		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)

		if err != nil || user.Token != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		// Check if the user's role is allowed
		allowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Your role does not have access."})
			c.Abort()
			return
		}

		// Attach claims to the context
		c.Set("userID", claims.UserID)

		// Proceed to the next handler
		c.Next()
	}
}

var TokenBlacklist = struct {
	sync.RWMutex
	tokens map[string]time.Time
}{
	tokens: make(map[string]time.Time),
}

// AddToBlacklist adds a token to the blacklist with an expiration time
func AddToBlacklist(token string, expiration time.Time) {
	TokenBlacklist.Lock()
	defer TokenBlacklist.Unlock()
	TokenBlacklist.tokens[token] = expiration
}

// IsBlacklisted checks if a token is in the blacklist
func IsBlacklisted(token string) bool {
	TokenBlacklist.RLock()
	defer TokenBlacklist.RUnlock()
	expiration, exists := TokenBlacklist.tokens[token]
	if !exists {
		return false
	}
	return time.Now().Before(expiration)
}
