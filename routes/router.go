package routes

import (
	"github.com/gin-gonic/gin"
	authRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/auth"
	cartRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/cart"
	productRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/product"
	productCategoryRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/product_category"
	userRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/user"
)

// RegisterRoutes registers all routes for the application
func RegisterRoutes(router *gin.Engine) {
	// Add other route groups here
	userRouter.RegisterUserRoutes(router)
	authRouter.RegisterAuthRoutes(router)
	productRouter.RegisterProductRoutes(router)
	productCategoryRouter.RegisterProductCategoryRoutes(router)
	cartRouter.RegisterCartRoutes(router)
}
