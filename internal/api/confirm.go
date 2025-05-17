package api

import (
	"net/http"

	"weatherApi/internal/model"
	"weatherApi/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

// confirmHandler validates the token and marks the subscription as confirmed
func confirmHandler(c *gin.Context) {
	token := c.Param("token")

	email, err := jwtutil.Parse(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	var sub model.Subscription
	if err := DB.Where("email = ?", email).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found / Subscription not found"})
		return
	}

	if sub.IsConfirmed {
		c.JSON(http.StatusOK, gin.H{"message": "Subscription already confirmed"})
		return
	}

	sub.IsConfirmed = true
	DB.Save(&sub)

	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}
