package routes

import (
	"github.com/gin-gonic/gin"
	authRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/auth"
	userRouter "github.com/wachirawittd123/shop-online-backend-golang/routes/user"
)

// RegisterRoutes registers all routes for the application
func RegisterRoutes(router *gin.Engine) {
	// Add other route groups here
	userRouter.RegisterUserRoutes(router)
	authRouter.RegisterAuthRoutes(router)
}
