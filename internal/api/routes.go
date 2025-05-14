package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/subscribe", subscribeHandler)
	api.GET("/confirm/:token", confirmHandler)
	api.GET("/unsubscribe/:token", unsubscribeHandler)
}
