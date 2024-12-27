package cartRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	cartController "github.com/wachirawittd123/shop-online-backend-golang/controller/cart"
)

// RegisterCartRoutes Routes defines cart-related routes
func RegisterCartRoutes(router *gin.Engine) {
	cartGroup := router.Group("/cart")
	{
		cartGroup.GET("/", common.AuthMiddleware("user", "admin"), cartController.GetCartWithProducts)
		cartGroup.GET("/:id", common.AuthMiddleware("user", "admin"), cartController.GetCartWithProducts)
		cartGroup.PUT("/", common.AuthMiddleware("user", "admin"), cartController.UpdateCart)
	}
}
