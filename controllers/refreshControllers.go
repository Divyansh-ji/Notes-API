package controllers

import (
	"main/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func RefreshToken(c *gin.Context) {
	// 1) Read refresh token from cookie
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(401, gin.H{"error": "No refresh token provided"})
		return
	}

	// 2) Parse and validate JWT
	claims, err := utils.ParseAndValidateJWT(rt)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// 3) Ensure it's a refresh token
	if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
		c.JSON(401, gin.H{"error": "Invalid token type"})
		return
	}

	// 4) Get user ID from claims
	uidFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}
	userID := uint(uidFloat)

	// 5) Issue new access token
	accessToken, err := utils.CreateJWT(userID, 15*time.Minute)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not create access token"})
		return
	}

	c.JSON(200, gin.H{"access_token": accessToken})
}
