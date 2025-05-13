package api

import (
	"net/http"
	"time"
	"weatherApi/internal/db"
	"weatherApi/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscribeRequest struct {
	Email     string `form:"email" binding:"required,email"`
	City      string `form:"city" binding:"required"`
	Frequency string `form:"frequency" binding:"required,oneof=daily hourly"`
}

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/subscribe", subscribeHandler)
}

func subscribeHandler(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Перевірка, чи вже є така підписка
	var existing model.Subscription
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already subscribed"})
		return
	}

	sub := model.Subscription{
		Email:       req.Email,
		City:        req.City,
		Frequency:   req.Frequency,
		IsConfirmed: false,
		Token:       uuid.New().String(),
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
		return
	}

	// TODO: надіслати лист із підтвердженням

	c.JSON(http.StatusOK, gin.H{"message": "Subscription received. Please confirm via email."})
}
