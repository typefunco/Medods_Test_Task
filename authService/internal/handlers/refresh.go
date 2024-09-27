package handlers

import (
	"log/slog"
	"medods_auth/authService/internal/entities"
	"medods_auth/authService/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) RefreshTokens(c *gin.Context) {
	refreshToken := c.Param("refresh_token")
	ip := c.ClientIP()

	storedToken, err := h.TokenRepo.GetRefreshTokenByToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if time.Now().After(storedToken.Exp) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}
	if storedToken.IP != ip {
		utils.SendWarningEmail(string(storedToken.UserID))
	}

	user, err := h.UserRepo.GetUserByID(storedToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	newAccessToken, err := h.JWT.GenerateToken(*user, ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	newRefreshToken, err := utils.GenerateRandomString(32)
	if err != nil {
		slog.Info("Can't generate Refresh Token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	newRefreshTokenHash, err := utils.HashToken(newRefreshToken)
	if err != nil {
		slog.Info("Can't generate hash Refresh Token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new refresh token"})
		return
	}

	newRefreshTokenEntity := entities.RefreshToken{
		TokenHash: newRefreshTokenHash,
		UserID:    user.ID,
		Exp:       time.Now().Add(24 * time.Hour * 168), // 1 week
		IP:        ip,
	}

	err = h.TokenRepo.SaveRefreshToken(newRefreshTokenEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new refresh token"})
		return
	}

	err = h.TokenRepo.DeleteRefreshToken(storedToken.TokenHash)
	if err != nil {
		slog.Info("Failed to delete old refresh token", "error", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
