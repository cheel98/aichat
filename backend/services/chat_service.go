package services

import (
	"aiChat/backend/database"
	"aiChat/backend/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// SaveChatSession 保存聊天会话
func SaveChatSession(session *models.ChatSession) error {
	db := database.GetDB()
	result := db.Create(session)
	return result.Error
}

// GetSessionsByUserID 获取用户的所有会话
func GetSessionsByUserID(userID uint) ([]models.ChatSession, error) {
	db := database.GetDB()
	var sessions []models.ChatSession

	// 使用GORM查询会话列表
	err := db.Where("user_id = ?", userID).
		Order("is_pinned DESC, updated_at DESC").
		Find(&sessions).Error

	if err != nil {
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}

	// 获取每个会话的消息数量和最后一条消息
	for i := range sessions {
		var messageCount int64
		db.Model(&models.ChatMessage{}).
			Where("session_id = ?", sessions[i].SessionID).
			Count(&messageCount)
		sessions[i].MessageCount = int(messageCount)

		var lastMessage models.ChatMessage
		err := db.Where("session_id = ?", sessions[i].SessionID).
			Order("created_at DESC").
			First(&lastMessage).Error
		if err == nil {
			sessions[i].LastMessage = &lastMessage
		}
	}

	return sessions, nil
}

// GetSessionByID 根据会话ID获取会话
func GetSessionByID(sessionID string) (*models.ChatSession, error) {
	db := database.GetDB()
	var session models.ChatSession

	err := db.Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("会话不存在")
		}
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}

	return &session, nil
}

// UpdateSession 更新会话信息
func UpdateSession(session *models.ChatSession) error {
	db := database.GetDB()
	result := db.Model(session).Updates(map[string]interface{}{
		"title":      session.Title,
		"is_pinned":  session.IsPinned,
		"updated_at": time.Now(),
	})
	return result.Error
}

// DeleteSession 删除会话及其消息
func DeleteSession(sessionID string) error {
	db := database.GetDB()
	return db.Transaction(func(tx *gorm.DB) error {
		// 删除会话的所有消息
		if err := tx.Where("session_id = ?", sessionID).Delete(&models.ChatMessage{}).Error; err != nil {
			return fmt.Errorf("删除会话消息失败: %w", err)
		}

		// 删除会话
		if err := tx.Where("session_id = ?", sessionID).Delete(&models.ChatSession{}).Error; err != nil {
			return fmt.Errorf("删除会话失败: %w", err)
		}

		return nil
	})
}

// SaveMessage 保存聊天消息
func SaveMessage(message *models.ChatMessage) error {
	db := database.GetDB()
	result := db.Create(message)
	return result.Error
}

// GetMessagesBySessionID 获取会话的所有消息
func GetMessagesBySessionID(sessionID string) ([]models.ChatMessage, error) {
	db := database.GetDB()
	var messages []models.ChatMessage

	err := db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}

	return messages, nil
}

// GetChatHistoryWithPagination 获取带分页的聊天历史
func GetChatHistoryWithPagination(sessionID string, page, pageSize int) ([]models.ChatMessage, int, error) {
	db := database.GetDB()
	var messages []models.ChatMessage
	var total int64

	// 获取消息总数
	if err := db.Model(&models.ChatMessage{}).
		Where("session_id = ?", sessionID).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取消息总数失败: %w", err)
	}

	// 获取分页消息
	offset := (page - 1) * pageSize
	err := db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&messages).Error

	if err != nil {
		return nil, 0, fmt.Errorf("查询消息失败: %w", err)
	}

	return messages, int(total), nil
}

// UpdateLastMessageTime 更新会话的最后消息时间
func UpdateLastMessageTime(sessionID string) error {
	db := database.GetDB()
	result := db.Model(&models.ChatSession{}).
		Where("session_id = ?", sessionID).
		Update("updated_at", time.Now())
	return result.Error
}
