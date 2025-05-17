package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires all API and UI routes.
// Only development-safe routes should be exposed in production builds.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/subscribe", subscribeHandler)
		api.GET("/confirm/:token", confirmHandler)
		api.GET("/unsubscribe/:token", unsubscribeHandler)

		// Debug/admin route — remove or protect in production
		api.GET("/subscriptions", listSubscriptionsHandler)

		api.GET("/weather", getWeatherHandler)
	}

	if gin.Mode() != gin.TestMode {
		r.LoadHTMLGlob("templates/*.html")
		r.Static("/static", "./static")
	}

	r.GET("/subscribe", func(c *gin.Context) {
		c.HTML(http.StatusOK, "subscribe.html", nil)
	})
}
