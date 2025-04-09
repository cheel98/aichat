package routes

import (
	"aiChat/backend/handlers"
	"aiChat/backend/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置路由
func SetupRoutes(r *gin.Engine) {
	// 全局中间件
	r.Use(middleware.CORSMiddleware())

	// 公共路由
	public := r.Group("/api")
	{
		// 认证相关
		auth := public.Group("/auth")
		{
			auth.POST("/register", handlers.RegisterHandler)
			auth.POST("/login", handlers.LoginHandler)
		}
	}

	// 需要认证的路由
	private := r.Group("/api")
	private.Use(middleware.AuthMiddleware())
	{
		// 用户相关
		user := private.Group("/user")
		{
			user.GET("/profile", handlers.GetUserProfileHandler)
			user.PUT("/profile", handlers.UpdateUserProfileHandler)
			user.PUT("/password", handlers.UpdatePasswordHandler)
			user.GET("/settings", handlers.GetUserSettingsHandler)
			user.PUT("/settings", handlers.UpdateUserSettingsHandler)
			user.POST("/logout", handlers.LogoutHandler)
		}
		// 聊天相关的需认证路由
		chat := private.Group("/chat")
		{
			chat.POST("/sessions", handlers.CreateSessionHandler)           // 创建会话
			chat.GET("/sessions", handlers.GetSessionsHandler)              // 获取会话列表
			chat.GET("/sessions/:id", handlers.GetSessionHandler)           // 获取会话详情
			chat.PUT("/sessions/:id", handlers.UpdateSessionHandler)        // 更新会话
			chat.DELETE("/sessions/:id", handlers.DeleteSessionHandler)     // 删除会话
			chat.POST("/sessions/:id", handlers.SendMessageHandler)         // 发送消息（流式返回）
			chat.POST("/retry", handlers.RetryMessageHandler)               // 重试生成回答
			chat.PUT("/response/active", handlers.SetActiveResponseHandler) // or /response/active 设置活跃回答
		}
	}
}
