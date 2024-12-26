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

func GetProductsCategory(c *gin.Context) {
	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	collection, ctx := common.GetCollection("product_category")

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product category"})
		return
	}
	defer cursor.Close(ctx)

	var productCategores []models.ProductCategory
	if err := cursor.All(ctx, &productCategores); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode product category"})
		return
	}

	if productCategores == nil {
		productCategores = []models.ProductCategory{}
	}

	c.JSON(http.StatusOK, gin.H{"product_category": productCategores, "status_code": http.StatusOK})
}

func AddProductCategory(c *gin.Context) {
	var productCategory models.ProductCategory
	if err := c.ShouldBindJSON(&productCategory); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if isCategoryNameTaken(productCategory.Name, c) {
		return
	}

	if err := insertCategory(productCategory, c); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product category added successfully", "status_code": http.StatusOK})
}

func UpdateProductCategory(c *gin.Context) {
	objectID, err := common.ConvertIDMongodb(c.Param("id"), c)
	if err != nil {
		return
	}

	var requestBody RequestUpdateBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	update := bson.M{"$set": bson.M{
		"name":       requestBody.Name,
		"updated_on": primitive.NewDateTimeFromTime(time.Now()),
	}}

	if err := common.UpdateOneCommonInDB(objectID, update, c, "product_category"); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product category updated successfully", "status_code": http.StatusOK})
}

func RemoveProductCategory(c *gin.Context) {
	objectID, err := common.ConvertIDMongodb(c.Param("id"), c)
	if err != nil {
		return
	}

	if isCategoryInUse(objectID, c) {
		c.JSON(http.StatusConflict, gin.H{"message": "Category is in use by existing products", "status_code": http.StatusConflict})
		return
	}

	if err := common.DeleteOneCommonByID(objectID, c, "product_category"); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully", "status_code": http.StatusOK})
}
