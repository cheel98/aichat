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
			auth.POST("/logout", handlers.LogoutHandler)
		}

		// 聊天相关的无需认证路由
		chat := public.Group("/chat")
		{
			chat.POST("/", handlers.ChatHandler) // 这个是之前的临时路由，用于兼容之前的前端
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
		}

		// 聊天相关的需要认证路由
		// chat := private.Group("/chat")
		{
			// 这些路由需要实现对应的handler
			// chat.POST("/sessions", handlers.CreateSessionHandler)
			// chat.GET("/sessions", handlers.GetSessionsHandler)
			// chat.GET("/sessions/:id", handlers.GetSessionHandler)
			// chat.PUT("/sessions/:id", handlers.UpdateSessionHandler)
			// chat.DELETE("/sessions/:id", handlers.DeleteSessionHandler)
			// chat.POST("/sessions/:id/messages", handlers.SendMessageHandler)
			// chat.GET("/sessions/:id/messages", handlers.GetMessagesHandler)
		}
	}
}
