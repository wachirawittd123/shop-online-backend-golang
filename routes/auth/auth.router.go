package authRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	authController "github.com/wachirawittd123/shop-online-backend-golang/controller/auth"
)

// RegisterUserRoutes defines user-related routes
func RegisterAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/logout", common.AuthMiddleware("user", "admin"), authController.Logout)
	}
}
