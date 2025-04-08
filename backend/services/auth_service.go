package services

import (
	"aiChat/backend/config"
	"aiChat/backend/database"
	"aiChat/backend/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrInvalidToken 无效令牌错误
	ErrInvalidToken = errors.New("invalid token")

	// ErrExpiredToken 过期令牌错误
	ErrExpiredToken = errors.New("token has expired")
)

// TokenClaims 令牌声明
type TokenClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint) (string, error) {
	db := database.GetDB()

	// 设置令牌过期时间
	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JWT.ExpiresIn) * time.Hour)

	// 创建JWT声明
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名字符串
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWT.Secret))
	if err != nil {
		return "", err
	}

	// 存储会话到数据库
	userSession := models.UserSession{
		UserID:     userID,
		Token:      tokenString,
		ExpireTime: expirationTime,
	}

	if err := db.Create(&userSession).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证令牌
func ValidateToken(tokenString string) (*TokenClaims, error) {
	db := database.GetDB()

	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 从令牌提取claims
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 检查令牌是否在数据库中
	var count int64
	err = db.Model(&models.UserSession{}).
		Where("token = ? AND expire_time > ?", tokenString, time.Now()).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// InvalidateToken 使令牌失效
func InvalidateToken(tokenString string) error {
	db := database.GetDB()
	return db.Where("token = ?", tokenString).Delete(&models.UserSession{}).Error
}

// CleanupExpiredTokens 清理过期会话
func CleanupExpiredTokens() error {
	db := database.GetDB()
	return db.Where("expire_time <= ?", time.Now()).Delete(&models.UserSession{}).Error
}
