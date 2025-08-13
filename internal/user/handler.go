package user

import (
	"bankapi/internal/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GET /users → Tüm kullanıcıları getirir
func GetUsers(c *gin.Context) {
	println("👥 Tüm kullanıcılar getiriliyor...")

	var users []User
	if err := db.DB.Find(&users).Error; err != nil {
		println("❌ Kullanıcılar getirilemedi:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcılar getirilemedi", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ Kullanıcılar bulundu, sayı:", len(users))

	// map to safe response
	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}
	c.JSON(http.StatusOK, responses)
}

// POST /users → Yeni kullanıcı ekler
func CreateUser(c *gin.Context) {
	println("👤 Yeni kullanıcı oluşturuluyor...")

	var req CreateUserRequest

	// JSON'dan gelen veriyi al
	if err := c.ShouldBindJSON(&req); err != nil {
		println("❌ Kullanıcı verisi parse edilemedi:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri", "message": "Kullanıcı bilgileri hatalı"})
		return
	}

	println("✅ Kullanıcı verisi alındı:", req.Username, req.Email)

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		println("❌ Şifre hash'lenemedi:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre işlenemedi", "message": "Teknik bir hata oluştu"})
		return
	}

	println("🔐 Şifre hash'lendi")

	user := User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         "user",
	}

	// Veritabanına kaydet (unique ihlallerini yakala)
	if err := db.DB.Create(&user).Error; err != nil {
		if isUniqueViolation(err) {
			println("❌ Kullanıcı zaten mevcut")
			c.JSON(http.StatusConflict, gin.H{"error": "Kullanıcı zaten mevcut", "message": "Bu kullanıcı adı veya e-posta zaten kayıtlı"})
			return
		}
		println("❌ Kullanıcı oluşturulamadı:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı oluşturulamadı", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ Kullanıcı başarıyla oluşturuldu, ID:", user.ID)
	c.JSON(http.StatusCreated, user.ToResponse())
}

// GET /users/:id → Tek kullanıcı getir
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	println("👤 Kullanıcı aranıyor, ID:", id)

	var u User
	if err := db.DB.First(&u, id).Error; err != nil {
		println("❌ Kullanıcı bulunamadı:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı", "message": "Belirtilen ID'ye sahip kullanıcı bulunamadı"})
		return
	}

	println("✅ Kullanıcı bulundu:", u.Username)
	c.JSON(http.StatusOK, u.ToResponse())
}

type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Role     *string `json:"role"`
}

// PUT /users/:id → Güncelle
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	println("✏️ Kullanıcı güncelleniyor, ID:", id)

	var u User
	if err := db.DB.First(&u, id).Error; err != nil {
		println("❌ Güncellenecek kullanıcı bulunamadı:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı", "message": "Belirtilen ID'ye sahip kullanıcı bulunamadı"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		println("❌ Güncelleme verisi parse edilemedi:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz veri", "message": "Güncelleme bilgileri hatalı"})
		return
	}

	println("✅ Güncelleme verisi alındı")

	if req.Username != nil {
		println("📝 Kullanıcı adı güncelleniyor:", *req.Username)
		u.Username = *req.Username
	}
	if req.Email != nil {
		println("📧 E-posta güncelleniyor:", *req.Email)
		u.Email = *req.Email
	}
	if req.Role != nil {
		println("👑 Rol güncelleniyor:", *req.Role)
		u.Role = *req.Role
	}

	if err := db.DB.Save(&u).Error; err != nil {
		if isUniqueViolation(err) {
			println("❌ Kullanıcı adı veya e-posta zaten kayıtlı")
			c.JSON(http.StatusConflict, gin.H{"error": "Kullanıcı zaten mevcut", "message": "Bu kullanıcı adı veya e-posta zaten kayıtlı"})
			return
		}
		println("❌ Kullanıcı güncellenemedi:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Güncellenemedi", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ Kullanıcı başarıyla güncellendi")
	c.JSON(http.StatusOK, u.ToResponse())
}

// DELETE /users/:id → Sil
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	println("🗑️ Kullanıcı siliniyor, ID:", id)

	if err := db.DB.Delete(&User{}, id).Error; err != nil {
		println("❌ Kullanıcı silinemedi:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Silinemedi", "message": "Teknik bir hata oluştu"})
		return
	}

	println("✅ Kullanıcı başarıyla silindi")
	c.Status(http.StatusNoContent)
}

// naive unique detection without pg driver types for simplicity
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	e := strings.ToLower(err.Error())
	return strings.Contains(e, "duplicate key") || strings.Contains(e, "unique constraint")
}
