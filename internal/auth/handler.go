package auth

import (
	"bankapi/internal/config"
	"bankapi/internal/db"
	"bankapi/internal/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func Register(c *gin.Context) {
	println("📝 Kullanıcı kaydı başlatılıyor...")

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("❌ Kayıt verisi parse edilemedi:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri", "message": "Kayıt bilgileri hatalı"})
		return
	}

	println("✅ Kayıt verisi alındı, kullanıcı:", req.Username, "email:", req.Email)

	// Şifreyi hash'le
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		println("❌ Şifre hash'lenemedi:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre işlenemedi", "message": "Teknik bir hata oluştu"})
		return
	}

	println("🔐 Şifre hash'lendi")

	// Kullanıcıyı oluştur
	newUser := user.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	// Veritabanına kaydet
	if err := db.DB.Create(&newUser).Error; err != nil {
		if isUniqueViolation(err) {
			println("❌ Kullanıcı zaten mevcut")
			c.JSON(http.StatusConflict, gin.H{"error": "Kullanıcı zaten mevcut", "message": "Bu kullanıcı adı veya e-posta zaten kayıtlı"})
			return
		}
		println("❌ Kullanıcı oluşturulamadı:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı oluşturulamadı", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ Kullanıcı başarıyla oluşturuldu, ID:", newUser.ID)
	c.JSON(http.StatusCreated, newUser.ToResponse())
}

func Login(c *gin.Context) {
	println("🔑 Kullanıcı girişi başlatılıyor...")

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("❌ Giriş verisi parse edilemedi:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri", "message": "Giriş bilgileri hatalı"})
		return
	}

	println("✅ Giriş verisi alındı, email:", req.Email)

	// Kullanıcıyı bul
	var u user.User
	if err := db.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		println("❌ Kullanıcı bulunamadı:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş başarısız", "message": "E-posta veya şifre hatalı"})
		return
	}

	println("👤 Kullanıcı bulundu:", u.Username)

	// Şifreyi kontrol et
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		println("❌ Şifre yanlış")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Giriş başarısız", "message": "E-posta veya şifre hatalı"})
		return
	}

	println("✅ Şifre doğrulandı, token oluşturuluyor...")

	// JWT token oluştur
	tokenPair, err := GenerateTokenPair(u.ID, u.Role, &config.Config{JWTSecret: "test-secret-key", TokenTTL: "1h"})
	if err != nil {
		println("❌ Token oluşturulamadı:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ JWT token başarıyla oluşturuldu")
	c.JSON(http.StatusOK, gin.H{
		"message":       "Giriş başarılı",
		"user":          u.ToResponse(),
		"token":         tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func Refresh(c *gin.Context) {
	println("🔄 Token yenileme başlatılıyor...")

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("❌ Refresh verisi parse edilemedi:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri", "message": "Refresh token bilgisi hatalı"})
		return
	}

	println("✅ Refresh token alındı")

	// Refresh token'ı doğrula
	token, err := ParseAndValidate(req.RefreshToken, &config.Config{JWTSecret: "test-secret-key"})
	if err != nil {
		println("❌ Refresh token geçersiz:", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz refresh token", "message": "Token yenilenemedi"})
		return
	}

	// Token'dan kullanıcı bilgilerini al
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		println("❌ Token claims parse edilemedi")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token", "message": "Token yenilenemedi"})
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		println("❌ User ID alınamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token", "message": "Token yenilenemedi"})
		return
	}

	role, ok := claims["role"].(string)
	if !ok {
		println("❌ Role alınamadı")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token", "message": "Token yenilenemedi"})
		return
	}

	// Yeni token çifti oluştur
	tokenPair, err := GenerateTokenPair(uint(userID), role, &config.Config{JWTSecret: "test-secret-key", TokenTTL: "1h"})
	if err != nil {
		println("❌ Yeni token oluşturulamadı:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulamadı", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ Token başarıyla yenilendi")
	c.JSON(http.StatusOK, gin.H{
		"message":       "Token yenilendi",
		"token":         tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

// naive unique detection without pg driver types for simplicity
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	e := strings.ToLower(err.Error())
	return strings.Contains(e, "duplicate key") || strings.Contains(e, "unique constraint")
}
