package services

import (
	"aiChat/backend/database"
	"aiChat/backend/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
		// 删除会话的所有AI响应
		if err := tx.Where("session_id = ?", sessionID).Delete(&models.AIResponse{}).Error; err != nil {
			return fmt.Errorf("删除会话AI响应失败: %w", err)
		}

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
	// 如果消息没有MessageID，则生成一个
	if message.MessageID == "" {
		message.MessageID = uuid.New().String()
	}

	db := database.GetDB()
	result := db.Create(message)
	return result.Error
}

// SaveAIResponse 保存AI响应
func SaveAIResponse(response *models.AIResponse) error {
	db := database.GetDB()
	result := db.Create(response)
	return result.Error
}

// GetMessagesBySessionID 获取会话的所有消息
func GetMessagesBySessionID(sessionID string) ([]models.ChatMessage, error) {
	db := database.GetDB()
	var messages []models.ChatMessage

	// 获取活跃的消息
	err := db.Where("session_id = ? AND is_active = ?", sessionID, true).
		Order("created_at ASC").
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}

	// 加载每条AI消息的替代回答
	for i, msg := range messages {
		if msg.Role == "ai" {
			alternatives, err := GetAIResponsesByMessageID(msg.MessageID)
			if err == nil && len(alternatives) > 0 {
				messages[i].AlternativeResponses = alternatives
			}
		}
	}

	return messages, nil
}

// GetMessageByID 根据消息ID获取消息
func GetMessageByID(messageID string) (*models.ChatMessage, error) {
	db := database.GetDB()
	var message models.ChatMessage

	err := db.Where("message_id = ?", messageID).First(&message).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("消息不存在")
		}
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}

	return &message, nil
}

// GetAIResponsesByMessageID 根据原始消息ID获取所有AI回答
func GetAIResponsesByMessageID(messageID string) ([]models.AIResponse, error) {
	db := database.GetDB()
	var responses []models.AIResponse

	err := db.Where("message_id = ?", messageID).
		Order("version ASC").
		Find(&responses).Error

	if err != nil {
		return nil, fmt.Errorf("查询AI回答失败: %w", err)
	}

	return responses, nil
}

// SetActiveAIResponse 设置活跃的AI响应
func SetActiveAIResponse(messageID string, version int) error {
	db := database.GetDB()

	// 开始事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 先将所有相关的AI响应和消息设为非活跃
		if err := tx.Model(&models.ChatMessage{}).
			Where("message_id = ?", messageID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.AIResponse{}).
			Where("message_id = ?", messageID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// 2. 根据版本号决定激活原始消息还是替代响应
		if version == 1 {
			// 激活原始消息
			if err := tx.Model(&models.ChatMessage{}).
				Where("message_id = ? AND version = ?", messageID, 1).
				Update("is_active", true).Error; err != nil {
				return err
			}
		} else {
			// 激活特定版本的替代响应
			if err := tx.Model(&models.AIResponse{}).
				Where("message_id = ? AND version = ?", messageID, version).
				Update("is_active", true).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetChatHistoryWithPagination 获取带分页的聊天历史
func GetChatHistoryWithPagination(sessionID string, page, pageSize int) ([]models.ChatMessage, int, error) {
	db := database.GetDB()
	var messages []models.ChatMessage
	var total int64

	// 获取消息总数
	if err := db.Model(&models.ChatMessage{}).
		Where("session_id = ? AND is_active = ?", sessionID, true).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取消息总数失败: %w", err)
	}

	// 获取分页消息
	offset := (page - 1) * pageSize
	err := db.Where("session_id = ? AND is_active = ?", sessionID, true).
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
