package models

import (
	"database/sql"
	"time"
)

// User 用户模型
type User struct {
	ID            uint64         `json:"id"`
	Username      string         `json:"username"`
	Password      string         `json:"-"` // 不序列化到JSON
	Email         sql.NullString `json:"email"`
	Phone         sql.NullString `json:"phone"`
	Avatar        sql.NullString `json:"avatar"`
	Status        int            `json:"status"`
	LoginType     int            `json:"login_type"`
	LastLoginTime sql.NullTime   `json:"last_login_time"`
	LastLoginIP   sql.NullString `json:"last_login_ip"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// UserSession 用户会话模型
type UserSession struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"user_id"`
	Token      string    `json:"token"`
	ExpireTime time.Time `json:"expire_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UserSettings 用户设置模型
type UserSettings struct {
	ID                  uint64    `json:"id"`
	UserID              uint64    `json:"user_id"`
	Theme               string    `json:"theme"`
	Language            string    `json:"language"`
	NotificationEnabled int       `json:"notification_enabled"`
	Prompt              string    `json:"prompt"`
	Rules               string    `json:"rules"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Password  string `json:"password" binding:"required,min=6,max=100"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,min=5,max=20"`
	LoginType int    `json:"login_type" binding:"required,oneof=1 2"` // 1: email, 2: phone
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account   string `json:"account" binding:"required"` // 邮箱或手机号
	Password  string `json:"password" binding:"required"`
	LoginType int    `json:"login_type" binding:"required,oneof=1 2"` // 1: email, 2: phone
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
