package productRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	productController "github.com/wachirawittd123/shop-online-backend-golang/controller/product"
)

// RegisterProductRoutes defines product-related routes
func RegisterProductRoutes(router *gin.Engine) {
	productGroup := router.Group("/products")
	{
		productGroup.GET("/", common.AuthMiddleware("user", "admin"), productController.GetProducts)
		productGroup.GET("/:id", common.AuthMiddleware("user", "admin"), productController.GetProduct)
		productGroup.POST("/", common.AuthMiddleware("admin"), productController.AddProduct)
		productGroup.PUT("/:id", common.AuthMiddleware("admin"), productController.UpdateProduct)
		productGroup.DELETE("/:id", common.AuthMiddleware("admin"), productController.RemoveProduct)
	}
}
