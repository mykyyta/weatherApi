package api

import (
	"net/http"

	"weatherApi/internal/model"

	"github.com/gin-gonic/gin"
)

func listSubscriptionsHandler(c *gin.Context) {
	var subs []model.Subscription

	if err := DB.Find(&subs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve subscriptions"})
		return
	}

	c.JSON(http.StatusOK, subs)
}
