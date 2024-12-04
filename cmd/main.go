package main

import (
	"github.com/Zhaoyikaiii/docmind/internal/api/controllers"
	"github.com/Zhaoyikaiii/docmind/internal/api/middleware"
	"github.com/Zhaoyikaiii/docmind/pkg/config"
	"github.com/Zhaoyikaiii/docmind/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/Zhaoyikaiii/docmind/internal/service"
	"github.com/Zhaoyikaiii/docmind/internal/interfaces"
	"github.com/Zhaoyikaiii/docmind/internal/storage"
	"github.com/Zhaoyikaiii/docmind/internal/storage/operators"
	"fmt"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	middleware.ApplyMiddleware(r)

	// Public routes
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Auth routes
	authController := &controllers.AuthController{}
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/refresh", authController.RefreshToken)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			username, _ := c.Get("username")
			c.JSON(200, gin.H{
				"message":  "You have access to protected resource",
				"user_id":  userID,
				"username": username,
			})
		})
	}

	return r
}

func setupFileOperator(cfg *config.Config) (interfaces.FileOperator, error) {
	switch cfg.Storage.Type {
	case "local":
		return storage.NewFileOperator(storage.Local, operators.LocalConfig{
			UploadDir:    cfg.Storage.LocalPath,
			MaxFileSize:  cfg.Storage.MaxFileSize,
			AllowedTypes: cfg.Storage.AllowedTypes,
		})

	case "oss":
		return storage.NewFileOperator(storage.OSS, operators.OSSConfig{
			Endpoint:        cfg.OSS.Endpoint,
			AccessKeyID:     cfg.OSS.AccessKeyID,
			AccessKeySecret: cfg.OSS.AccessKeySecret,
				BucketName:      cfg.OSS.BucketName,
				BasePath:        cfg.OSS.BasePath,
				MaxFileSize:     cfg.OSS.MaxFileSize,
				AllowedTypes:    cfg.OSS.AllowedTypes,
		})

	// ... 其他存储类型的配置
	}

	return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
}

func main() {
	config.LoadConfig()

	utils.InitLogger()
	defer utils.Logger.Sync()

	r := setupRouter()

	port := config.GetString("server.port")
	utils.Logger.Info("Starting server on port " + port)

	if err := r.Run(":" + port); err != nil {
		utils.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
