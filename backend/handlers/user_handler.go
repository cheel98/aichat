package handlers

import (
	"errors"
	"net/http"
	"strings"

	"aiChat/backend/models"
	"aiChat/backend/services"

	"github.com/gin-gonic/gin"
)

// RegisterHandler 用户注册处理器
func RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 验证登录类型
	if req.LoginType != 1 && req.LoginType != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的登录类型"})
		return
	}

	// 验证邮箱或手机号
	if req.LoginType == 1 && req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱登录类型需要提供邮箱"})
		return
	}

	if req.LoginType == 2 && req.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号登录类型需要提供手机号"})
		return
	}

	// 注册用户
	userID, err := services.RegisterUser(req)
	if err != nil {
		if errors.Is(err, services.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		} else if errors.Is(err, services.ErrEmailExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "邮箱已存在"})
		} else if errors.Is(err, services.ErrPhoneExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "手机号已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user_id": userID,
	})
}

// LoginHandler 用户登录处理器
func LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 登录
	resp, err := services.LoginUser(req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		} else if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "账号或密码错误"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// LogoutHandler 用户登出处理器
func LogoutHandler(c *gin.Context) {
	// 从请求头获取Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少认证令牌"})
		return
	}

	// 提取令牌
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的认证格式"})
		return
	}

	// 使令牌失效
	if err := services.InvalidateToken(tokenString); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

// GetUserProfileHandler 获取用户资料处理器
func GetUserProfileHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
		return
	}

	user, err := services.GetUserByID(userID.(uint))
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserProfileHandler 更新用户资料处理器
func UpdateUserProfileHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 更新用户资料
	err := services.UpdateUserProfile(userID.(uint), req)
	if err != nil {
		if errors.Is(err, services.ErrUsernameExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "用户名已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新资料失败: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "资料更新成功"})
}

// UpdatePasswordHandler 更新密码处理器
func UpdatePasswordHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
		return
	}

	var req models.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 更新密码
	err := services.UpdateUserPassword(userID.(uint64), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "原密码错误"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码更新成功"})
}

// GetUserSettingsHandler 获取用户设置处理器
func GetUserSettingsHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
		return
	}

	settings, err := services.GetUserSettings(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取设置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateUserSettingsHandler 更新用户设置处理器
func UpdateUserSettingsHandler(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权，请先登录"})
		return
	}

	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 更新设置
	err := services.UpdateUserSettings(userID.(uint64), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新设置失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设置更新成功"})
}
