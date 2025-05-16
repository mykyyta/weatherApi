package api

import (
	"fmt"
	"net/http"
	"time"

	"weatherApi/internal/model"
	emailutil "weatherApi/pkg/email"
	"weatherApi/pkg/jwtutil"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type SubscribeRequest struct {
	Email     string `form:"email" binding:"required,email"`
	City      string `form:"city" binding:"required"`
	Frequency string `form:"frequency" binding:"required,oneof=daily hourly"`
}

func subscribeHandler(c *gin.Context) {
	var req SubscribeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existing model.Subscription
	err := DB.Where("email = ?", req.Email).First(&existing).Error

	if err == nil {
		if existing.IsConfirmed && !existing.IsUnsubscribed {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already subscribed"})
			return
		}
		_ = DB.Delete(&existing)
	}

	token, err := jwtutil.Generate(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	sub := model.Subscription{
		ID:             uuid.New().String(),
		Email:          req.Email,
		City:           req.City,
		Frequency:      req.Frequency,
		IsConfirmed:    false,
		IsUnsubscribed: false,
		Token:          token,
		CreatedAt:      time.Now(),
	}

	if err := DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save subscription"})
		return
	}

	confirmURL := fmt.Sprintf("http://localhost:8080/api/confirm/%s", token)
	fmt.Println("Confirm URL:", confirmURL)

	go func() {
		if err := emailutil.SendConfirmationEmail(req.Email, token); err != nil {
			fmt.Printf("⚠️ Failed to send confirmation email to %s: %v\n", req.Email, err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Subscription successful. Confirmation email sent."})
}

func confirmHandler(c *gin.Context) {
	token := c.Param("token")

	email, err := jwtutil.Parse(token)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid or expired token"})
		return
	}

	var sub model.Subscription
	if err := DB.Where("email = ?", email).First(&sub).Error; err != nil {
		c.JSON(404, gin.H{"error": "not found"})
		return
	}

	sub.IsConfirmed = true
	DB.Save(&sub)

	c.JSON(200, gin.H{"message": "subscription confirmed"})
}
