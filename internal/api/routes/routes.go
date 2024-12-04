package routes

import (
	"github.com/Zhaoyikaiii/docmind/internal/api/controllers"
	"github.com/Zhaoyikaiii/docmind/internal/api/handlers"
	"github.com/Zhaoyikaiii/docmind/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, dc *controllers.DocumentController, uh *handlers.UploadHandler, fc *controllers.FileController) {
	// Apply global middleware
	middleware.ApplyMiddleware(r)

	// Public routes
	public := r.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			authController := controllers.AuthController{}
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
		}
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// Document routes
		docs := protected.Group("/documents")
		{
			docs.POST("", dc.CreateDocument)
			docs.PUT("/:id", dc.UpdateDocument)
			docs.DELETE("/:id", dc.DeleteDocument)
			docs.GET("/:id", dc.GetDocument)
			docs.GET("", dc.ListDocuments)
			docs.POST("/:id/versions", dc.CreateVersion)
			docs.GET("/:id/versions", dc.GetVersions)
			docs.POST("/:id/tags", dc.ManageTags)
		}

		// File upload routes
		upload := protected.Group("/upload")
		{
			upload.POST("/file", uh.HandleFileUpload)
		}

		// File routes
		files := protected.Group("/files")
		{
			files.GET("/:id", fc.GetFile)
			files.GET("", fc.ListFiles)
			files.DELETE("/:id", fc.DeleteFile)
			files.POST("/:id/document", fc.AssociateWithDocument)
		}
	}
} 