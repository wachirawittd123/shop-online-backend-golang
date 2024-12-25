package userRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/wachirawittd123/shop-online-backend-golang/common"
	userController "github.com/wachirawittd123/shop-online-backend-golang/controller/user"
)

// RegisterUserRoutes defines user-related routes
func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.GET("/", common.AuthMiddleware("admin"), userController.GetUsers)
		userGroup.POST("/", common.AuthMiddleware("admin"), userController.AddUser)
		userGroup.DELETE("/:id", common.AuthMiddleware("admin"), userController.RemoveUser)
	}
}
