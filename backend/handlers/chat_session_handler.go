package handlers

import (
	"aiChat/backend/database"
	"aiChat/backend/models"
	"aiChat/backend/services"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateSessionHandler 创建新的聊天会话
func CreateSessionHandler(c *gin.Context) {
	var req models.CreateSessionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取当前用户ID（从认证中间件中设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 创建新会话
	session := models.ChatSession{
		SessionID: uuid.New().String(), // 生成唯一会话ID
		UserID:    userID.(uint),
		Title:     req.Title,
		IsPinned:  0,
	}

	// 调用服务保存会话
	if err := services.SaveChatSession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建会话失败"})
		return
	}

	// 返回响应，包含会话ID
	c.JSON(http.StatusCreated, gin.H{
		"id":         session.ID,
		"session_id": session.SessionID,
		"title":      session.Title,
		"created_at": session.CreatedAt,
	})
}

// GetSessionsHandler 获取当前用户的所有聊天会话
func GetSessionsHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 调用服务获取会话列表
	sessions, err := services.GetSessionsByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"conversations": sessions,
		"total":         len(sessions),
	})
}

// GetSessionHandler 获取特定会话的详情
func GetSessionHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 调用服务获取会话详情
	session, err := services.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话详情失败"})
		return
	}

	// 验证会话所有权
	if session.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此会话"})
		return
	}

	// 获取会话的消息列表
	messages, err := services.GetMessagesBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": session.SessionID,
		"title":      session.Title,
		"messages":   messages,
	})
}

// UpdateSessionHandler 更新会话信息
func UpdateSessionHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	var req models.UpdateSessionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取现有会话
	session, err := services.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话失败"})
		return
	}

	// 验证会话所有权
	if session.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此会话"})
		return
	}

	// 更新会话信息
	if req.Title != "" {
		session.Title = req.Title
	}

	if req.IsPinned != nil {
		if *req.IsPinned {
			session.IsPinned = 1
		} else {
			session.IsPinned = 0
		}
	}

	// 保存更新
	if err := services.UpdateSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新会话失败"})
		return
	}

	c.JSON(http.StatusOK, session)
}

// DeleteSessionHandler 删除会话
func DeleteSessionHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取现有会话
	session, err := services.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话失败"})
		return
	}

	// 验证会话所有权
	if session.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此会话"})
		return
	}

	// 删除会话及相关消息
	if err := services.DeleteSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "会话已删除"})
}

// SendMessageHandler 在指定会话中发送消息，并流式返回AI响应
func SendMessageHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	var req struct {
		Content      string `json:"content" binding:"required"`
		DeepThinking bool   `json:"thinking"`             // 是否使用深度思考模式（使用reasoner模型）
		MessageID    string `json:"message_id,omitempty"` // 可选的消息ID，用于重试
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式有误"})
		return
	}
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息内容不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 验证会话存在且属于当前用户
	session, err := services.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话失败"})
		return
	}

	if session.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此会话"})
		return
	}

	// 生成消息ID
	messageID := req.MessageID
	if messageID == "" {
		messageID = uuid.New().String()
	}

	// 判断是否是重试
	isRetry := req.MessageID != ""
	var version int = 1

	if isRetry {
		// 如果是重试，先检查原始消息是否存在
		_, err := services.GetMessageByID(messageID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无法找到原始消息"})
			return
		}

		// 获取所有响应，确定新的版本号
		responses, err := services.GetAIResponsesByMessageID(messageID)
		if err == nil {
			version = len(responses) + 2 // 原始消息是版本1，新回答版本递增
		} else {
			version = 2 // 没有找到额外回答，这是第一次重试
		}
	} else {
		// 创建用户消息
		userMessage := models.ChatMessage{
			UserID:    userID.(uint),
			SessionID: sessionID,
			Role:      "user",
			MessageID: messageID,
			Content:   req.Content,
			IsActive:  true,
			Version:   1,
		}

		// 保存用户消息
		if err := services.SaveMessage(&userMessage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存消息失败"})
			return
		}
	}

	// 设置流式响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // 禁用Nginx缓冲

	// 刷新响应头
	c.Writer.Flush()

	// 添加一个包装器，在每次写入后刷新缓冲区
	flushWriter := &FlushWriter{Writer: c.Writer}

	// 获取AI服务实例
	aiService := services.GetDefaultDeepSeekService()

	// 调用流式API获取回复，传递deep_thinking参数
	aiReply, aiThinking, err := aiService.StreamChatResponse(req.Content, flushWriter, req.DeepThinking)
	if err != nil {
		c.Error(err)
		return
	}

	// 根据是否是重试，决定保存到哪个表
	if isRetry {
		// 创建新的AI响应记录
		aiResponse := models.AIResponse{
			MessageID:    messageID,
			SessionID:    sessionID,
			ThinkContent: aiThinking,
			Content:      aiReply,
			Version:      version,
			IsActive:     false, // 默认不激活新回答
		}

		// 保存AI响应
		if err := services.SaveAIResponse(&aiResponse); err != nil {
			fmt.Printf("保存AI回复失败: %v\n", err)
		}

		// 返回响应版本信息
		c.Writer.Write([]byte(fmt.Sprintf("\n\n$responseVersion$%d", version)))
	} else {
		// 创建AI消息并保存到数据库（使用完整的回复内容）
		aiMessage := models.ChatMessage{
			UserID:       0, // AI消息的UserID为0
			SessionID:    sessionID,
			Role:         "ai",
			MessageID:    messageID,
			ThinkContent: aiThinking,
			Content:      aiReply,
			Version:      1,
			IsActive:     true,
		}

		// 保存AI消息
		if err := services.SaveMessage(&aiMessage); err != nil {
			// 仅记录错误，不中断响应
			fmt.Printf("保存AI回复失败: %v\n", err)
		}

		// 返回消息ID
		c.Writer.Write([]byte(fmt.Sprintf("\n\n$messageId$%s", messageID)))
	}

	// 更新会话的最后更新时间
	if err := services.UpdateSession(session); err != nil {
		// 仅记录错误，不中断响应
		fmt.Printf("更新会话失败: %v\n", err)
	}
}

// FlushWriter 用于确保每次写入后立即刷新
type FlushWriter struct {
	Writer gin.ResponseWriter
}

// Write 实现io.Writer接口，在写入后立即刷新
func (fw *FlushWriter) Write(data []byte) (int, error) {
	n, err := fw.Writer.Write(data)
	fw.Writer.Flush()
	return n, err
}

// GetMessagesHandler 获取指定会话的消息列表
func GetMessagesHandler(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 验证会话存在且属于当前用户
	session, err := services.GetSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话失败"})
		return
	}

	if session.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此会话"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20 // 默认每页20条消息
	}

	// 获取消息列表
	messages, err := services.GetMessagesBySessionID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    len(messages),
	})
}

// RetryMessageHandler 重试消息，生成新的AI回答
func RetryMessageHandler(c *gin.Context) {
	var req struct {
		MessageID    string `json:"message_id" binding:"required"`
		DeepThinking bool   `json:"thinking"` // 是否使用深度思考模式
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式有误"})
		return
	}

	// 获取当前用户ID
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取原始消息
	message, err := services.GetMessageByID(req.MessageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "找不到原始消息"})
		return
	}

	var userMessage models.ChatMessage
	db := database.GetDB()

	// 判断传入的是用户消息还是AI消息
	if message.Role == "user" {
		// 如果传入的是用户消息ID，直接使用
		userMessage = *message
	} else if message.Role == "ai" {
		// 如果传入的是AI消息ID，查找对应的用户消息
		err = db.Where("session_id = ? AND role = ? AND created_at < ?",
			message.SessionID, "user", message.CreatedAt).
			Order("created_at DESC").
			First(&userMessage).Error

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "找不到用户问题"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息角色无效"})
		return
	}

	// 构建重试请求
	retryReq := struct {
		Content      string `json:"content"`
		DeepThinking bool   `json:"thinking"`
		MessageID    string `json:"message_id"`
	}{
		Content:      userMessage.Content,
		DeepThinking: req.DeepThinking,
		MessageID:    req.MessageID,
	}

	// 将原始请求参数传递给SendMessageHandler
	c.Set("retryRequest", retryReq)

	// 设置参数并调用SendMessageHandler
	c.Params = append(c.Params, gin.Param{Key: "id", Value: message.SessionID})
	c.Request.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(`{"content":"%s","thinking":%v,"message_id":"%s"}`,
		userMessage.Content, req.DeepThinking, req.MessageID)))

	SendMessageHandler(c)
}

// SetActiveResponseHandler 设置当前活跃的AI回答版本
func SetActiveResponseHandler(c *gin.Context) {
	var req struct {
		MessageID string `json:"message_id" binding:"required"`
		Version   int    `json:"version" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式有误"})
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
		return
	}

	// 获取消息
	message, err := services.GetMessageByID(req.MessageID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "找不到消息"})
		return
	}

	// 验证用户是否有权限修改此会话
	session, err := services.GetSessionByID(message.SessionID)
	if err != nil || session.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此会话"})
		return
	}

	// 设置活跃版本
	if err := services.SetActiveAIResponse(req.MessageID, req.Version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设置活跃版本失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已设置活跃版本"})
}
