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

func GetProducts(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	collection, ctx := common.GetCollection("products")

	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"name": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

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

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
		return
	}

	if products == nil {
		products = []models.Product{}
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func GetProduct(c *gin.Context) {
	_id := c.Param("id")

	collection, ctx := common.GetCollection("products")

	product := GetProductById(_id, collection, ctx, c)

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func AddProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection, ctx := common.GetCollection("products")

	var existingProduct models.Product
	err := collection.FindOne(ctx, bson.M{"name": product.Name}).Decode(&existingProduct)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product name already in use"})
		return
	}

	var dateNow = primitive.NewDateTimeFromTime(time.Now())
	product.ID = primitive.NewObjectID()
	product.CreatedAt = dateNow
	product.UpdatedOn = dateNow

	_, err = collection.InsertOne(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully"})
}

func UpdateProduct(c *gin.Context) {
	_id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	var requestBody struct {
		Name   string `json:"name" binding:"required"`
		Price  int    `json:"price" binding:"required"`
		Detail string `json:"detail"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	collection, ctx := common.GetCollection("products")
	update := bson.M{"$set": bson.M{
		"name":       requestBody.Name,
		"price":      requestBody.Price,
		"detail":     requestBody.Detail,
		"updated_on": primitive.NewDateTimeFromTime(time.Now()),
	}}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product name"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func RemoveProduct(c *gin.Context) {
	_id := c.Param("id")

	objectID, err := common.ConvertIDMongodb(_id, c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
		return
	}

	collection, ctx := common.GetCollection("products")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
