package productController

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProducts(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	collection, ctx := common.GetCollection("products")

	// Build filters for the aggregation pipeline
	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
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

	// Define the aggregation pipeline
	pipeline := []bson.M{
		{"$match": filter}, // Apply the filters
		{
			"$lookup": bson.M{
				"from":         "product_category", // The collection to join
				"localField":   "id_category",      // Field in the "products" collection
				"foreignField": "_id",              // Field in the "product_category" collection
				"as":           "category_details", // Output array field
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$category_details", // Unwind the joined array
				"preserveNullAndEmptyArrays": false,               // Exclude products with no category match
			},
		},
	}

	// Execute the aggregation pipeline
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch products", "status_code": http.StatusInternalServerError})
		return
	}
	defer cursor.Close(ctx)

	// Decode the result into a slice of maps (for dynamic fields)
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode products", "status_code": http.StatusInternalServerError})
		return
	}

	if results == nil {
		results = []bson.M{}
	}

	// Respond with the joined data
	c.JSON(http.StatusOK, gin.H{"products": results, "status_code": http.StatusOK})
}

func GetProduct(c *gin.Context) {
	_id := c.Param("id")

	collection, ctx := common.GetCollection("products")

	product := GetProductById(_id, collection, ctx, c)

	c.JSON(http.StatusOK, gin.H{"product": product, "status_code": http.StatusOK})
}

func AddProduct(c *gin.Context) {
	// Parse and validate request body
	var requestBody AddProductRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate category ID
	idCategory, err := common.ConvertIDMongodb(requestBody.IDCategory, c)
	if err != nil {
		return
	}

	// Check if product name already exists
	if isProductNameTaken(requestBody.Name, c) {
		return
	}

	// Prepare the product for insertion
	product := prepareProductForInsertion(requestBody, idCategory)

	// Insert the product into the database
	if err := insertProduct(product, c); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added successfully", "status_code": http.StatusOK})
}

func UpdateProduct(c *gin.Context) {
	// Parse and validate the product ID
	objectID, err := common.ConvertIDMongodb(c.Param("id"), c)
	if err != nil {
		return
	}

	// Parse and validate the request body
	var requestBody AddProductRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Convert category ID to ObjectID
	idCategory, err := primitive.ObjectIDFromHex(requestBody.IDCategory)
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid category ID format", err)
		return
	}

	// Prepare the update document
	update := bson.M{
		"$set": bson.M{
			"name":        requestBody.Name,
			"price":       requestBody.Price,
			"detail":      requestBody.Detail,
			"id_category": idCategory,
			"updated_on":  primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	// Perform the update operation
	if err := common.UpdateOneCommonInDB(objectID, update, c, "products"); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully", "status_code": http.StatusOK})
}

func RemoveProduct(c *gin.Context) {
	// Parse and validate the product ID
	objectID, err := common.ConvertIDMongodb(c.Param("id"), c)
	if err != nil {
		return
	}

	// Attempt to delete the product
	if err := common.DeleteOneCommonByID(objectID, c, "products"); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully", "status_code": http.StatusOK})
}
