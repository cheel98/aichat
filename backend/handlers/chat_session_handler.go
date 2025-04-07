package handlers

import (
	"aiChat/backend/models"
	"aiChat/backend/services"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 生成随机会话ID
func generateSessionID() string {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 生成16位的随机字符串
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := make([]byte, 16)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}

	// 添加时间戳前缀
	timestamp := time.Now().Unix()
	return fmt.Sprintf("session_%d_%s", timestamp, string(id))
}

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
		UserID:    userID.(uint64),
		Title:     req.Title,
		IsPinned:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
	sessions, err := services.GetSessionsByUserID(userID.(uint64))
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
	if session.UserID != userID.(uint64) {
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
		"id":       session.ID,
		"title":    session.Title,
		"messages": messages,
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
	if session.UserID != userID.(uint64) {
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

	session.UpdatedAt = time.Now()

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
	if session.UserID != userID.(uint64) {
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
		Content string `json:"content" binding:"required"`
		Message string `json:"message"` // 兼容旧版API
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	// 兼容处理，如果没有content但有message字段，则使用message
	content := req.Content
	if content == "" && req.Message != "" {
		content = req.Message
	}

	if content == "" {
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

	if session.UserID != userID.(uint64) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问此会话"})
		return
	}

	// 创建用户消息
	userMessage := models.ChatMessage{
		UserID:    userID.(uint64),
		SessionID: sessionID,
		Role:      "user",
		Content:   content,
		CreatedAt: time.Now(),
	}

	// 保存用户消息
	if err := services.SaveMessage(&userMessage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存消息失败"})
		return
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

	// 调用流式API获取回复
	aiReply, err := aiService.StreamChatResponse(content, flushWriter)
	if err != nil {
		c.Error(err)
		return
	}

	// 创建AI消息并保存到数据库（使用完整的回复内容）
	aiMessage := models.ChatMessage{
		UserID:    0, // AI消息的UserID为0
		SessionID: sessionID,
		Role:      "ai",
		Content:   aiReply,
		CreatedAt: time.Now(),
	}

	// 保存AI消息
	if err := services.SaveMessage(&aiMessage); err != nil {
		// 仅记录错误，不中断响应
		fmt.Printf("保存AI回复失败: %v\n", err)
	}

	// 更新会话的最后更新时间
	session.UpdatedAt = time.Now()
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

	if session.UserID != userID.(uint64) {
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
