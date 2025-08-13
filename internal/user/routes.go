package user

import (
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, middlewares ...gin.HandlerFunc) {
	userRoutes := router.Group("/api/v1/users", middlewares...)
	{
		userRoutes.GET("/", GetUsers)
		userRoutes.POST("/", CreateUser)
		userRoutes.GET(":id", GetUserByID)
		userRoutes.PUT(":id", UpdateUser)
		userRoutes.DELETE(":id", DeleteUser)
	}
}
