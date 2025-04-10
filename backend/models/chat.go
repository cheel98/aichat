package models

import (
	"time"

	"gorm.io/gorm"
)

// ChatSession 聊天会话模型
type ChatSession struct {
	gorm.Model
	SessionID string `gorm:"type:varchar(50);not null;unique" json:"session_id"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	Title     string `gorm:"type:varchar(100);not null" json:"title"`
	IsPinned  int    `gorm:"type:tinyint;default:0" json:"is_pinned"`
	// 这些字段不存储在数据库中，用于前端显示
	MessageCount int          `gorm:"-" json:"message_count,omitempty"`
	LastMessage  *ChatMessage `gorm:"-" json:"last_message,omitempty"`
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	gorm.Model
	UserID       uint   `gorm:"not null" json:"user_id"`
	SessionID    string `gorm:"type:varchar(50);not null" json:"session_id"`
	Role         string `gorm:"type:varchar(20);not null" json:"role"`       // user 或 ai
	MessageID    string `gorm:"type:varchar(50);not null" json:"message_id"` // 消息唯一ID
	ThinkContent string `gorm:"type:text;not null" json:"think_content"`
	Content      string `gorm:"type:text;not null" json:"content"`
	Version      int    `gorm:"type:int;default:1" json:"version"`          // 版本号，用于跟踪重试
	IsActive     bool   `gorm:"type:boolean;default:true" json:"is_active"` // 当前是否是活跃版本
	// 关联的其他响应
	AlternativeResponses []AIResponse `gorm:"-" json:"alternative_responses,omitempty"`
}

// AIResponse AI响应模型，用于存储同一问题的多个回答
type AIResponse struct {
	gorm.Model
	MessageID    string `gorm:"type:varchar(50);not null" json:"message_id"` // 关联到原始消息
	SessionID    string `gorm:"type:varchar(50);not null" json:"session_id"`
	ThinkContent string `gorm:"type:text;not null" json:"think_content"`
	Content      string `gorm:"type:text;not null" json:"content"`
	Version      int    `gorm:"type:int;not null" json:"version"`            // 版本号
	IsActive     bool   `gorm:"type:boolean;default:false" json:"is_active"` // 是否是当前活跃版本
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Title string `json:"title"`
}

// UpdateSessionRequest 更新会话请求
type UpdateSessionRequest struct {
	Title    string `json:"title"`
	IsPinned *bool  `json:"is_pinned,omitempty"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// RetryMessageRequest 重试消息请求
type RetryMessageRequest struct {
	MessageID string `json:"message_id" binding:"required"`
}

// AIMessageResponse AI回复消息
type AIMessageResponse struct {
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// SessionListResponse 会话列表响应
type SessionListResponse struct {
	Sessions []ChatSession `json:"sessions"`
	Total    int           `json:"total"`
}

// MessageListResponse 消息列表响应
type MessageListResponse struct {
	Messages []ChatMessage `json:"messages"`
	Total    int           `json:"total"`
}
