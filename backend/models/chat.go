package models

import (
	"time"
)

// ChatSession 聊天会话模型
type ChatSession struct {
	ID        uint64    `json:"id"`
	SessionID string    `json:"session_id"`
	UserID    uint64    `json:"user_id"`
	Title     string    `json:"title"`
	IsPinned  int       `json:"is_pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// 这些字段不存储在数据库中，用于前端显示
	MessageCount int          `json:"message_count,omitempty"`
	LastMessage  *ChatMessage `json:"last_message,omitempty"`
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"` // user 或 ai
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
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
