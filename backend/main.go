package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aiChat/backend/config" // 请根据实际项目结构调整导入路径
	"aiChat/backend/database"
	"aiChat/backend/routes"
	"aiChat/backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	appConfig, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化DeepSeek服务
	deepseekService := services.GetDefaultDeepSeekService()
	deepseekService.SetAPIKey(appConfig.DeepSeek.APIKey)
	deepseekService.SetModel(appConfig.DeepSeek.Model)

	// 初始化数据库连接
	if err := database.InitDB(); err != nil {
		log.Fatalf("初始化数据库连接失败: %v", err)
	}
	defer database.CloseDB()

	// 确保数据库表结构
	if err := database.SetupDatabase(); err != nil {
		log.Fatalf("设置数据库表结构失败: %v", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 设置路由
	routes.SetupRoutes(r)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port),
		Handler: r,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("服务器启动在 %s:%d", appConfig.Server.Host, appConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("监听错误: %v", err)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 是 syscall.SIGINT
	// kill -9 是 syscall.SIGKILL，但是无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("关闭服务器...")

	// 设置5秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器优雅退出")
}
