package balance

import (
	"github.com/gin-gonic/gin"
)

func RegisterBalanceRoutes(router *gin.Engine) {
	RegisterRoutes(router)
}
