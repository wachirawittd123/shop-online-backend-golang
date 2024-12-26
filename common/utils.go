package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertIDMongodb(_id string, c *gin.Context) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid category ID format", err)
		return primitive.NilObjectID, err
	}
	return objectID, nil
}

// respondWithError sends an error response
func RespondWithError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		c.JSON(statusCode, gin.H{"message": message, "details": err.Error(), "status_code": statusCode})
	} else {
		c.JSON(statusCode, gin.H{"message": message, "status_code": statusCode})
	}
}

// update performs the update operation in the database
func UpdateOneCommonInDB(objectID primitive.ObjectID, update bson.M, c *gin.Context, collectionDB string) error {
	collection, ctx := GetCollection(collectionDB)
	defer ctx.Done()

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to update product", err)
		return err
	}

	if result.ModifiedCount == 0 {
		RespondWithError(c, http.StatusNotFound, "Product not found", nil)
		return err
	}

	return nil
}

// deletes a product by its ObjectID from the database
func DeleteOneCommonByID(objectID primitive.ObjectID, c *gin.Context, collectionDB string) error {
	collection, ctx := GetCollection(collectionDB)
	defer ctx.Done()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, "Failed to delete "+collectionDB, err)
		return err
	}

	if result.DeletedCount == 0 {
		RespondWithError(c, http.StatusNotFound, collectionDB+" not found", nil)
		return err
	}

	return nil
}
