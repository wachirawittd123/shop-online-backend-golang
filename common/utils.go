package common

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertIDMongodb(_id string, c *gin.Context) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		log.Println("Invalid UserID format:", _id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return primitive.NilObjectID, nil
	}
	return objectID, nil
}
