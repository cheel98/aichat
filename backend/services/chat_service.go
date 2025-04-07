package services

import (
	"aiChat/backend/database"
	"aiChat/backend/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// SaveChatSession 保存聊天会话
func SaveChatSession(session *models.ChatSession) error {
	db := database.GetDB()

	// 插入会话记录
	result, err := db.Exec(
		"INSERT INTO chat_sessions (session_id, user_id, title, is_pinned, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		session.SessionID, session.UserID, session.Title, session.IsPinned, session.CreatedAt, session.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存会话失败: %w", err)
	}

	// 获取自动生成的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取会话ID失败: %w", err)
	}

	session.ID = uint64(id)
	return nil
}

// GetSessionsByUserID 获取用户的所有会话
func GetSessionsByUserID(userID uint64) ([]models.ChatSession, error) {
	db := database.GetDB()

	// 查询用户的所有会话
	rows, err := db.Query(`
		SELECT s.id, s.session_id, s.user_id, s.title, s.is_pinned, s.created_at, s.updated_at,
			   COUNT(m.id) as message_count,
			   (
				   SELECT JSON_OBJECT(
					   'id', m2.id,
					   'role', m2.role,
					   'content', SUBSTRING(m2.content, 1, 100),
					   'created_at', m2.created_at
				   )
				   FROM chat_messages m2
				   WHERE m2.session_id = s.session_id
				   ORDER BY m2.created_at DESC
				   LIMIT 1
			   ) as last_message
		FROM chat_sessions s
		LEFT JOIN chat_messages m ON s.session_id = m.session_id
		WHERE s.user_id = ?
		GROUP BY s.id
		ORDER BY s.is_pinned DESC, s.updated_at DESC
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}
	defer rows.Close()

	// 解析结果
	var sessions []models.ChatSession
	for rows.Next() {
		var session models.ChatSession
		var messageCount int
		var lastMessageJSON sql.NullString

		err := rows.Scan(
			&session.ID,
			&session.SessionID,
			&session.UserID,
			&session.Title,
			&session.IsPinned,
			&session.CreatedAt,
			&session.UpdatedAt,
			&messageCount,
			&lastMessageJSON,
		)

		if err != nil {
			return nil, fmt.Errorf("解析会话数据失败: %w", err)
		}

		session.MessageCount = messageCount

		// 处理最后一条消息的JSON
		if lastMessageJSON.Valid {
			// 这里应该使用json.Unmarshal来解析JSON
			// 简化处理，实际项目中应该正确解析JSON到LastMessage字段
			// var lastMessage models.ChatMessage
			// json.Unmarshal([]byte(lastMessageJSON.String), &lastMessage)
			// session.LastMessage = &lastMessage
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetSessionByID 根据会话ID获取会话
func GetSessionByID(sessionID string) (*models.ChatSession, error) {
	db := database.GetDB()

	var session models.ChatSession
	err := db.QueryRow(
		"SELECT id, session_id, user_id, title, is_pinned, created_at, updated_at FROM chat_sessions WHERE session_id = ?",
		sessionID,
	).Scan(
		&session.ID,
		&session.SessionID,
		&session.UserID,
		&session.Title,
		&session.IsPinned,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("会话不存在")
		}
		return nil, fmt.Errorf("查询会话失败: %w", err)
	}

	return &session, nil
}

// UpdateSession 更新会话信息
func UpdateSession(session *models.ChatSession) error {
	db := database.GetDB()

	_, err := db.Exec(
		"UPDATE chat_sessions SET title = ?, is_pinned = ?, updated_at = ? WHERE session_id = ?",
		session.Title, session.IsPinned, session.UpdatedAt, session.SessionID,
	)

	if err != nil {
		return fmt.Errorf("更新会话失败: %w", err)
	}

	return nil
}

// DeleteSession 删除会话及其消息
func DeleteSession(sessionID string) error {
	db := database.GetDB()

	// 启动事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}

	// 删除会话的所有消息
	_, err = tx.Exec("DELETE FROM chat_messages WHERE session_id = ?", sessionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("删除会话消息失败: %w", err)
	}

	// 删除会话
	_, err = tx.Exec("DELETE FROM chat_sessions WHERE session_id = ?", sessionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("删除会话失败: %w", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}

// SaveMessage 保存聊天消息
func SaveMessage(message *models.ChatMessage) error {
	db := database.GetDB()

	result, err := db.Exec(
		"INSERT INTO chat_messages (user_id, session_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)",
		message.UserID, message.SessionID, message.Role, message.Content, message.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("保存消息失败: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取消息ID失败: %w", err)
	}

	message.ID = uint64(id)
	return nil
}

// GetMessagesBySessionID 获取会话的所有消息
func GetMessagesBySessionID(sessionID string) ([]models.ChatMessage, error) {
	db := database.GetDB()

	rows, err := db.Query(
		"SELECT id, user_id, session_id, role, content, created_at FROM chat_messages WHERE session_id = ? ORDER BY created_at ASC",
		sessionID,
	)

	if err != nil {
		return nil, fmt.Errorf("查询消息失败: %w", err)
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage

		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("解析消息数据失败: %w", err)
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// GetChatHistoryWithPagination 获取带分页的聊天历史
func GetChatHistoryWithPagination(sessionID string, page, pageSize int) ([]models.ChatMessage, int, error) {
	db := database.GetDB()

	// 获取消息总数
	var total int
	err := db.QueryRow("SELECT COUNT(*) FROM chat_messages WHERE session_id = ?", sessionID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("获取消息总数失败: %w", err)
	}

	// 计算分页
	offset := (page - 1) * pageSize

	// 获取分页消息
	rows, err := db.Query(
		"SELECT id, user_id, session_id, role, content, created_at FROM chat_messages WHERE session_id = ? ORDER BY created_at ASC LIMIT ? OFFSET ?",
		sessionID, pageSize, offset,
	)

	if err != nil {
		return nil, 0, fmt.Errorf("查询消息失败: %w", err)
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage

		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("解析消息数据失败: %w", err)
		}

		messages = append(messages, message)
	}

	return messages, total, nil
}

// UpdateLastMessageTime 更新会话的最后消息时间
func UpdateLastMessageTime(sessionID string) error {
	db := database.GetDB()

	_, err := db.Exec(
		"UPDATE chat_sessions SET updated_at = ? WHERE session_id = ?",
		time.Now(), sessionID,
	)

	if err != nil {
		return fmt.Errorf("更新会话时间失败: %w", err)
	}

	return nil
}
