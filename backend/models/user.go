package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username      string     `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password      string     `gorm:"type:varchar(255);not null" json:"-"`
	Email         string     `gorm:"type:varchar(100);unique" json:"email"`
	Phone         string     `gorm:"type:varchar(20);unique" json:"phone"`
	Avatar        string     `gorm:"type:varchar(255)" json:"avatar"`
	Status        int        `gorm:"type:tinyint;default:1" json:"status"`
	LoginType     int        `gorm:"type:tinyint;not null" json:"login_type"`
	LastLoginTime *time.Time `json:"last_login_time"`
	LastLoginIP   string     `gorm:"type:varchar(50)" json:"last_login_ip"`
}

// UserSession 用户会话模型
type UserSession struct {
	gorm.Model
	UserID     uint      `gorm:"not null" json:"user_id"`
	Token      string    `gorm:"type:varchar(255);not null;unique" json:"token"`
	ExpireTime time.Time `gorm:"not null" json:"expire_time"`
}

// UserSettings 用户设置模型
type UserSettings struct {
	gorm.Model
	UserID              uint   `gorm:"not null;unique" json:"user_id"`
	Theme               string `gorm:"type:varchar(20);default:'light'" json:"theme"`
	Language            string `gorm:"type:varchar(10);default:'zh-CN'" json:"language"`
	NotificationEnabled int    `gorm:"type:tinyint;default:1" json:"notification_enabled"`
	Prompt              string `gorm:"type:text" json:"prompt"`
	Rules               string `gorm:"type:text" json:"rules"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6,max=100"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,min=5,max=20"`
	LoginType int    `json:"login_type" binding:"required,oneof=1 2"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account   string `json:"account" binding:"required"`
	Password  string `json:"password" binding:"required"`
	LoginType int    `json:"login_type" binding:"required,oneof=1 2"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Avatar   string `json:"avatar"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	Theme               string `json:"theme" binding:"omitempty,oneof=dark light"`
	Language            string `json:"language" binding:"omitempty"`
	NotificationEnabled bool   `json:"notification_enabled"`
	Prompt              string `json:"prompt" binding:"omitempty"`
	Rules               string `json:"rules" binding:"omitempty"`
}
