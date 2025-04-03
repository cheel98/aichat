# AI聊天项目

这是一个使用Vue和Golang构建的AI聊天应用，用于对接DeepSeek AI。

## 项目结构

项目采用前后端分离的架构：

```
.
├── frontend/  # Vue前端项目
└── backend/   # Golang后端项目
```

## 后端（Golang）

后端提供API接口，用于将用户消息转发给DeepSeek AI并返回应答。

### 配置

在运行前需要配置`backend/config.yml`文件，特别是添加您的DeepSeek API密钥：

```yaml
deepseek:
  api_key: your-deepseek-api-key-here # 替换为您的API密钥
  base_url: https://api.deepseek.com/v1
  model: deepseek-v3 # 默认使用v3模型
```

### 安装和运行

```bash
# 进入后端目录
cd backend

# 运行后端服务
go run main.go
```

后端服务将在 http://localhost:8080 上运行。

## 前端（Vue）

前端提供用户界面，用于发送消息和显示AI回复。

### 安装和运行

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run serve
```

前端服务将在 http://localhost:8081 上运行。

## 基本功能

当前实现的基本功能是：
- 用户可以发送消息到后端
- 后端使用DeepSeek API（默认v3模型）处理用户消息并返回智能回复
- 支持DeepSeek API的配置（API密钥、基础URL和模型选择）

## 后续开发计划

- 添加用户认证
- 添加消息历史记录
- 添加对话上下文管理
- 支持更多DeepSeek模型和参数配置
- 添加流式响应支持 