package handlers

import (
	"github.com/gin-gonic/gin"
)

type AuthHandlers interface {
	GenerateTokens(c *gin.Context)
	RefreshTokens(c *gin.Context)
}

func RegisterRoutes(router *gin.Engine, authHandle AuthHandlers) {
	auth := router.Group("/auth")
	{
		auth.POST("/login/:user_id", authHandle.GenerateTokens)
		auth.POST("/refresh/:refresh_token", authHandle.RefreshTokens)
	}
}
