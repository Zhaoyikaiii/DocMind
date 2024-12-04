package controllers

import (
	"github.com/Zhaoyikaiii/docmind/pkg/auth"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthController handles authentication related requests
type AuthController struct {
	// 这里可以添加用户服务等依赖
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: 验证用户凭据
	// 这里应该调用用户服务来验证用户名和密码
	// 现在我们先使用模拟数据
	userID := uint(1)
	username := req.Username
	role := "user"

	tokenInfo, err := auth.GenerateToken(userID, username, role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(200, tokenInfo)
}

// RefreshToken handles token refresh
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	tokenInfo, err := auth.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(200, tokenInfo)
}
