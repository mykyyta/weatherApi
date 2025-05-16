package api

import (
	"net/http"

	"weatherApi/pkg/weatherapi"

	"github.com/gin-gonic/gin"
)

func getWeatherHandler(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City is required"})
		return
	}

	weather, statusCode, err := weatherapi.FetchWithStatus(city)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, weather)
}
