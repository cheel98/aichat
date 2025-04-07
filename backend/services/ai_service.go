package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
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
	Stream      bool          `json:"stream"`
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

// StreamChatResponse 流式获取聊天回复，同时返回完整的响应文本
func (s *DeepSeekService) StreamChatResponse(userMessage string, writer io.Writer) (string, error) {
	if s.Config.APIKey == "" {
		response := "未配置API密钥，无法连接DeepSeek服务"
		writer.Write([]byte(response))
		return response, nil
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
		Stream:      true, // 启用流式响应
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
	req.Header.Set("Accept", "text/event-stream")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取响应流
	reader := bufio.NewReader(resp.Body)

	// 用于收集完整回复
	var fullResponse strings.Builder

	// 如果DeepSeek API不支持流式输出或出现问题，我们实现一个模拟的流式响应
	// 这里通过简单地逐字符输出来模拟流式效果
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		// 读取整个响应
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("读取响应失败: %v", err)
		}

		var response ChatResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return "", fmt.Errorf("解析响应失败: %v", err)
		}

		if len(response.Choices) == 0 {
			return "", fmt.Errorf("没有收到有效回复")
		}

		// 获取完整回复
		fullReply := response.Choices[0].Message.Content

		// 模拟流式输出，每次发送一个字符
		for _, char := range fullReply {
			_, err := writer.Write([]byte(string(char)))
			if err != nil {
				return fullResponse.String(), fmt.Errorf("写入响应失败: %v", err)
			}
			// 小延迟，模拟打字效果
			time.Sleep(10 * time.Millisecond)
			fullResponse.WriteRune(char)
		}
		return fullResponse.String(), nil
	}

	// 处理真正的SSE流
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fullResponse.String(), fmt.Errorf("读取流失败: %v", err)
		}

		// 跳过空行
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 处理SSE格式数据
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// 检查流是否结束
			if data == "[DONE]" {
				break
			}

			// 解析JSON
			var streamResponse struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
					FinishReason *string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &streamResponse); err != nil {
				continue // 忽略解析错误，继续处理
			}

			// 检查是否有内容需要写入
			if len(streamResponse.Choices) > 0 {
				content := streamResponse.Choices[0].Delta.Content
				if content != "" {
					_, err := writer.Write([]byte(content))
					if err != nil {
						return fullResponse.String(), fmt.Errorf("写入响应失败: %v", err)
					}
					// 收集完整回复
					fullResponse.WriteString(content)
				}

				// 检查是否完成
				if streamResponse.Choices[0].FinishReason != nil {
					break
				}
			}
		}
	}

	return fullResponse.String(), nil
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
