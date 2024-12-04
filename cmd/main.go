package main

import (
	"github.com/Zhaoyikaiii/docmind/internal/api/middleware"
	"github.com/Zhaoyikaiii/docmind/pkg/config"
	"github.com/Zhaoyikaiii/docmind/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter() *gin.Engine {
	r := gin.New()

	middleware.ApplyMiddleware(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":    "pong",
			"request_id": c.GetString("RequestID"),
		})
	})

	authorized := r.Group("/api")
	authorized.Use(middleware.NewMiddleware().Auth)
	{
		authorized.GET("/protected", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message":    "You have access to protected resource",
				"request_id": c.GetString("RequestID"),
			})
		})
	}

	return r
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
