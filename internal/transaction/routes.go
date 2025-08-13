package transaction

import (
	"github.com/gin-gonic/gin"
)

func RegisterTransactionRoutes(router *gin.Engine) {
	RegisterRoutes(router)
}
