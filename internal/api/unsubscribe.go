package api

import (
	"net/http"
	"weatherApi/internal/db"
	"weatherApi/internal/model"
	"weatherApi/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

func unsubscribeHandler(c *gin.Context) {
	token := c.Param("token")

	email, err := jwtutil.Parse(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	var sub model.Subscription
	if err := db.DB.Where("email = ?", email).First(&sub).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	if sub.IsUnsubscribed {
		c.JSON(http.StatusOK, gin.H{"message": "You are already unsubscribed"})
		return
	}

	sub.IsUnsubscribed = true
	if err := db.DB.Save(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}
