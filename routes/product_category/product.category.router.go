package productCategoryRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	productCategoryController "github.com/wachirawittd123/shop-online-backend-golang/controller/product_category"
)

// RegisterProductCategory Routes defines product-category-related routes
func RegisterProductCategoryRoutes(router *gin.Engine) {
	productCategoryGroup := router.Group("/product-category")
	{
		productCategoryGroup.GET("/", common.AuthMiddleware("user", "admin"), productCategoryController.GetProductsCategory)
		productCategoryGroup.POST("/", common.AuthMiddleware("admin"), productCategoryController.AddProductCategory)
		productCategoryGroup.PUT("/:id", common.AuthMiddleware("admin"), productCategoryController.UpdateProductCategory)
		productCategoryGroup.DELETE("/:id", common.AuthMiddleware("admin"), productCategoryController.RemoveProductCategory)
	}
}
