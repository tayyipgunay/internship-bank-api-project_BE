package auth

import "github.com/gin-gonic/gin"

func RegisterAuthRoutes(router *gin.Engine) {
	r := router.Group("/api/v1/auth")
	{
		r.POST("/register", Register)
		r.POST("/login", Login)
		r.POST("/refresh", Refresh)
	}
}
