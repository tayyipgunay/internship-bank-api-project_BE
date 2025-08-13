package middleware

import (
	"bankapi/internal/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}
	return func(c *gin.Context) {
		val, exists := c.Get(ContextUserKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz"})
			return
		}
		u, ok := val.(user.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz"})
			return
		}
		if _, ok := roleSet[u.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Erişim reddedildi"})
			return
		}
		c.Next()
	}
}
