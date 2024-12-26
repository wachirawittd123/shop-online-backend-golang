package productController

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	models "github.com/wachirawittd123/shop-online-backend-golang/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProductById(id string, collection *mongo.Collection, ctx context.Context, c *gin.Context) models.Product {
	var existingProduct = models.Product{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return existingProduct
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&existingProduct)
	if err != nil {
		// Product not found
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return existingProduct
	}
	return existingProduct
}

// isProductNameTaken checks if the product name is already in use
func isProductNameTaken(name string, c *gin.Context) bool {
	collection, ctx := common.GetCollection("products")

	var existingProduct models.Product
	err := collection.FindOne(ctx, bson.M{"name": name}).Decode(&existingProduct)
	if err == nil {
		common.RespondWithError(c, http.StatusConflict, "Product name already in use", nil)
		return true
	}
	return false
}

// prepareProductForInsertion prepares a product document for insertion
func prepareProductForInsertion(request AddProductRequest, idCategory primitive.ObjectID) models.Product {
	now := primitive.NewDateTimeFromTime(time.Now())
	return models.Product{
		ID:         primitive.NewObjectID(),
		Name:       request.Name,
		Price:      request.Price,
		Detail:     request.Detail,
		IDCategory: idCategory,
		CreatedAt:  now,
		UpdatedOn:  now,
	}
}

// insertProduct inserts a product into the database
func insertProduct(product models.Product, c *gin.Context) error {
	collection, ctx := common.GetCollection("products")

	_, err := collection.InsertOne(ctx, product)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to add product", err)
	}
	return err
}
