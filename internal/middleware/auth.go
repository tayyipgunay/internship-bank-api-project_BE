package middleware

import (
	"bankapi/internal/auth"
	"bankapi/internal/config"
	"bankapi/internal/db"
	"bankapi/internal/user"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserKey = "currentUser"
)

func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}
	return parts[1], nil
}

// AuthMiddleware creates authentication middleware
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return AuthRequired(cfg)
}

func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		println("🔐 Auth middleware çalışıyor - path:", c.Request.URL.Path)

		tokenStr, err := getBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			println("❌ Token alınamadı:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Yetkisiz erişim", "message": "Token bulunamadı"})
			return
		}

		println("🔑 Token alındı, doğrulanıyor...")

		t, err := auth.ParseAndValidate(tokenStr, cfg)
		if err != nil || !t.Valid {
			println("❌ Token geçersiz:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token", "message": "Token doğrulanamadı"})
			return
		}

		claims, ok := t.Claims.(jwt.MapClaims)
		if !ok {
			println("❌ Token claims okunamadı")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token hatası", "message": "Token bilgileri okunamadı"})
			return
		}

		subVal, ok := claims["sub"].(float64)
		if !ok {
			println("❌ Token subject geçersiz, tip:", fmt.Sprintf("%T", claims["sub"]))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token hatası", "message": "Token kullanıcı bilgisi geçersiz"})
			return
		}
		userID := uint(subVal)

		println("👤 Kullanıcı ID:", userID, "veritabanından aranıyor...")

		var u user.User
		if err := db.DB.First(&u, userID).Error; err != nil {
			println("❌ Kullanıcı bulunamadı:", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı bulunamadı", "message": "Token ile ilişkili kullanıcı bulunamadı"})
			return
		}

		println("✅ Kullanıcı doğrulandı:", u.Username, "rol:", u.Role)
		c.Set(ContextUserKey, u)
		c.Next()
	}
}
