package services

import (
	"aiChat/backend/config"
	"aiChat/backend/database"
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
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint64) (string, error) {
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
	_, err = database.DB.Exec(
		"INSERT INTO user_sessions (user_id, token, expire_time) VALUES (?, ?, ?)",
		userID, tokenString, expirationTime,
	)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证令牌
func ValidateToken(tokenString string) (*TokenClaims, error) {
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
	var count int
	err = database.DB.QueryRow(
		"SELECT COUNT(*) FROM user_sessions WHERE token = ? AND expire_time > NOW()",
		tokenString,
	).Scan(&count)

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
	_, err := database.DB.Exec("DELETE FROM user_sessions WHERE token = ?", tokenString)
	return err
}

// CleanupExpiredTokens 清理过期会话
func CleanupExpiredTokens() error {
	_, err := database.DB.Exec("DELETE FROM user_sessions WHERE expire_time <= NOW()")
	return err
}
