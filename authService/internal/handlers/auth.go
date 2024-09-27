package handlers

import (
	"log/slog"
	"medods_auth/authService/internal/entities"
	"medods_auth/authService/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserRepo  UserRepository
	TokenRepo TokenRepository
	JWT       utils.JWTService
}

type UserRepository interface {
	GetUserByID(id int) (*entities.User, error)
}

type TokenRepository interface {
	SaveRefreshToken(token entities.RefreshToken) error
	GetRefreshTokenByToken(tokenHash string) (*entities.RefreshToken, error)
	DeleteRefreshToken(tokenHash string) error
}

func (h *AuthHandler) GenerateTokens(c *gin.Context) {
	userID := c.Param("user_id")
	ip := c.ClientIP()

	usrID, err := strconv.Atoi(userID)
	if err != nil {
		slog.Error("Invalid user ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.UserRepo.GetUserByID(usrID)
	if err != nil {
		slog.Error("User not found", "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	accessToken, err := h.JWT.GenerateToken(*user, ip)
	if err != nil {
		slog.Error("Failed to generate access token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := utils.GenerateRandomString(32)
	if err != nil {
		slog.Error("Failed to generate refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	refreshTokenHash, err := utils.HashToken(refreshToken)
	if err != nil {
		slog.Error("Failed to hash refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process refresh token"})
		return
	}

	refreshTokenEntity := entities.RefreshToken{
		TokenHash: refreshTokenHash,
		UserID:    user.ID,
		Exp:       time.Now().Add(24 * time.Hour),
		IP:        ip,
	}

	err = h.TokenRepo.SaveRefreshToken(refreshTokenEntity)
	if err != nil {
		slog.Error("Failed to save refresh token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
