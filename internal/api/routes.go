package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// API routing
	api := r.Group("/api")
	{
		api.POST("/subscribe", subscribeHandler)
		api.GET("/confirm/:token", confirmHandler)
		api.GET("/unsubscribe/:token", unsubscribeHandler)
		api.GET("/subscriptions", listSubscriptionsHandler)
		api.GET("/weather", getWeatherHandler)
	}

	// HTML templates and static files (Bootstrap, styles, js, etc)
	if gin.Mode() != gin.TestMode {
		r.LoadHTMLGlob("templates/*.html")
		r.Static("/static", "./static")
	}

	// Subscription page
	r.GET("/subscribe", func(c *gin.Context) {
		c.HTML(http.StatusOK, "subscribe.html", nil)
	})
}
