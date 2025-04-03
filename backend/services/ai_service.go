package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DeepSeekConfig 保存DeepSeek API的配置信息
type DeepSeekConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

// DeepSeekService 处理与DeepSeek API的交互
type DeepSeekService struct {
	Config DeepSeekConfig
}

// 创建默认DeepSeek服务实例
var defaultDeepSeekService = NewDeepSeekService(DeepSeekConfig{
	BaseURL: "https://api.deepseek.com/v1",
	Model:   "deepseek-chat", // 默认使用v3模型
})

// NewDeepSeekService 创建一个新的DeepSeekService实例
func NewDeepSeekService(config DeepSeekConfig) *DeepSeekService {
	return &DeepSeekService{
		Config: config,
	}
}

// GetDefaultDeepSeekService 获取默认DeepSeek服务实例
func GetDefaultDeepSeekService() *DeepSeekService {
	return defaultDeepSeekService
}

// SetAPIKey 设置API密钥
func (s *DeepSeekService) SetAPIKey(apiKey string) {
	s.Config.APIKey = apiKey
}

// SetModel 设置要使用的模型
func (s *DeepSeekService) SetModel(model string) {
	s.Config.Model = model
}

// ChatRequest DeepSeek API聊天请求
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse DeepSeek API聊天响应
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      ChatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GetChatResponse 获取聊天回复
func (s *DeepSeekService) GetChatResponse(userMessage string) (string, error) {
	if s.Config.APIKey == "" {
		return "未配置API密钥，无法连接DeepSeek服务", nil
	}

	// 构建请求体
	requestBody := ChatRequest{
		Model: s.Config.Model,
		Messages: []ChatMessage{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
		Temperature: 0.7,
		MaxTokens:   2000,
	}

	// 转换为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/chat/completions", s.Config.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.Config.APIKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查是否有回复
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("没有收到有效回复")
	}

	// 返回AI的回复
	return response.Choices[0].Message.Content, nil
}
