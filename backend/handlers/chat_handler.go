package handlers

import (
	"net/http"

	"aiChat/backend/services" // 请根据实际项目结构调整导入路径

	"github.com/gin-gonic/gin"
)

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Message string `json:"message"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	Reply string `json:"reply"`
}

// ChatHandler 处理聊天请求
func ChatHandler(c *gin.Context) {
	var msg ChatMessage
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取DeepSeek服务实例
	aiService := services.GetDefaultDeepSeekService()

	// 调用DeepSeek API获取回复
	reply, err := aiService.GetChatResponse(msg.Message, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取AI回复失败: " + err.Error()})
		return
	}

	// 返回AI回复
	response := ChatResponse{
		Reply: reply,
	}

	c.JSON(http.StatusOK, response)
}
