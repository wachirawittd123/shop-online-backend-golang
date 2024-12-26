package productCategoryController

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// isCategoryNameTaken checks if the category name is already in use
func isCategoryNameTaken(name string, c *gin.Context) bool {
	collection, ctx := common.GetCollection("product_category")
	defer ctx.Done()

	var existingCategory models.ProductCategory
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&existingCategory)
	if err == nil {
		common.RespondWithError(c, http.StatusConflict, "Product category name already in use", nil)
		return true
	}
	return false
}

// insertCategory prepares and inserts the product category into the database
func insertCategory(productCategory models.ProductCategory, c *gin.Context) error {
	collection, ctx := common.GetCollection("product_category")
	defer ctx.Done()

	now := primitive.NewDateTimeFromTime(time.Now())
	productCategory.ID = primitive.NewObjectID()
	productCategory.CreatedAt = now
	productCategory.UpdatedOn = now

	_, err := collection.InsertOne(ctx, productCategory)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to add product category", err)
		return err
	}
	return nil
}

func isCategoryInUse(categoryID primitive.ObjectID, c *gin.Context) bool {
	collection, ctx := common.GetCollection("products")
	defer ctx.Done()

	// Check if any product references the category
	count, err := collection.CountDocuments(ctx, bson.M{"id_category": categoryID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category usage", "details": err.Error()})
		return true // Assume it's in use if we can't verify
	}

	return count > 0
}
