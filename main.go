package main

import (
	"ai-demo/config"
	"ai-demo/data"
	"ai-demo/db"
	"ai-demo/handler"
	"ai-demo/middleware"
	"ai-demo/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database, err := db.NewDB(cfg.DBDSN)
	if err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	if err := database.AutoMigrate(&data.User{}, &data.Bill{}); err != nil {
		log.Fatalf("database migrate failed: %v", err)
	}

	aiService := service.NewOpenAIService(cfg.APIKey, cfg.BaseURL, cfg.Model)
	weChatService := service.NewWeChatService(cfg.WeChatAppID, cfg.WeChatSecret, cfg.WeChatMockMode)

	authHandler := handler.NewAuthHandler(database, weChatService, cfg.JWTSecret)
	chatHandler := handler.NewChatHandler(aiService, database)

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/auth/login", authHandler.Login)

	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	authGroup.POST("/chat", chatHandler.Chat)
	authGroup.GET("/bills/today", chatHandler.TodayBill)
	authGroup.GET("/bills/history", chatHandler.HistoryBills)
	authGroup.PUT("/bills/:id", chatHandler.UpdateBill)
	authGroup.DELETE("/bills/:id", chatHandler.DeleteBill)

	r.Run(":" + cfg.Port)
}
