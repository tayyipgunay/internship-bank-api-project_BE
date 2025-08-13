package auth

import (
	"errors"
	"time"

	"bankapi/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateTokenPair(userID uint, role string, cfg *config.Config) (TokenPair, error) {
	println("🔑 Token çifti oluşturuluyor, kullanıcı ID:", userID, "rol:", role)

	if cfg.JWTSecret == "" {
		println("❌ JWT_SECRET konfigürasyonu eksik!")
		return TokenPair{}, errors.New("JWT secret is not configured")
	}

	ttl := time.Hour
	if cfg.TokenTTL != "" {
		if d, err := time.ParseDuration(cfg.TokenTTL); err == nil {
			ttl = d
			println("⏰ Token TTL:", ttl.String())
		} else {
			println("⚠️ Geçersiz Token TTL, varsayılan 1h kullanılıyor")
		}
	}

	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
		"typ":  "access",
	}
	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh",
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		println("❌ Access token oluşturulamadı:", err.Error())
		return TokenPair{}, err
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		println("❌ Refresh token oluşturulamadı:", err.Error())
		return TokenPair{}, err
	}

	println("✅ Token çifti başarıyla oluşturuldu")
	return TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func ParseAndValidate(tokenStr string, cfg *config.Config) (*jwt.Token, error) {
	println("🔍 Token doğrulanıyor...")

	if tokenStr == "" {
		println("❌ Token string boş")
		return nil, errors.New("token string is empty")
	}

	if cfg.JWTSecret == "" {
		println("❌ JWT_SECRET konfigürasyonu eksik")
		return nil, errors.New("JWT secret is not configured")
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			println("❌ Geçersiz imzalama metodu:", t.Method.Alg())
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		println("❌ Token parse hatası:", err.Error())
		return nil, err
	}

	if !token.Valid {
		println("❌ Token geçersiz")
		return nil, errors.New("token is invalid")
	}

	println("✅ Token başarıyla doğrulandı")
	return token, nil
}
