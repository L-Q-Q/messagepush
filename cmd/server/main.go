package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"message-push-system/config"
	"message-push-system/internal/handler"
	"message-push-system/internal/middleware"
	"message-push-system/internal/repository"
	"message-push-system/internal/service"
	"message-push-system/internal/smtp"
	"message-push-system/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	db, err := cfg.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Database connected successfully")

	// 初始化 Repository
	groupRepo := repository.NewGroupRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	logRepo := repository.NewLogRepository(db)

	// 初始化 Service
	groupService := service.NewGroupService(groupRepo)
	memberService := service.NewMemberService(memberRepo)
	messageService := service.NewMessageService(messageRepo, groupRepo)
	logService := service.NewLogService(logRepo)

	// 初始化 SMTP 客户端
	smtpClient := smtp.NewSMTPClient()

	// 初始化 Push Worker
	pushWorker := worker.NewPushWorker(smtpClient, messageRepo, memberRepo, logRepo, groupRepo)
	pushWorker.Start()
	defer pushWorker.Stop()

	// 初始化 Handler
	groupHandler := handler.NewGroupHandler(groupService)
	memberHandler := handler.NewMemberHandler(memberService)
	messageHandler := handler.NewMessageHandler(messageService, pushWorker)
	logHandler := handler.NewLogHandler(logService)

	// 初始化 Gin 路由
	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 添加错误处理中间件
	router.Use(middleware.ErrorHandler())

	// 静态文件服务
	router.Static("/static", "./web/static")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/index.html")
	})

	// API 路由组
	api := router.Group("/api/v1")
	{
		// 群组路由
		api.POST("/groups", groupHandler.Create)
		api.DELETE("/groups/:id", groupHandler.Delete)
		api.GET("/groups", groupHandler.List)

		// 成员路由
		api.POST("/groups/:group_id/members", memberHandler.Add)
		api.DELETE("/members/:id", memberHandler.Remove)
		api.GET("/groups/:group_id/members", memberHandler.ListByGroup)

		// 消息路由
		api.POST("/messages", messageHandler.Create)

		// 日志路由
		api.GET("/logs", logHandler.List)
	}

	// 启动 HTTP 服务器
	port := cfg.Server.Port
	log.Printf("Starting server on port %s", port)

	// 优雅关闭
	go func() {
		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
