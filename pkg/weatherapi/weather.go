package weatherapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"weatherApi/internal/model"
)

type weatherAPIResponse struct {
	Current struct {
		TempC     float64 `json:"temp_c"`
		Humidity  int     `json:"humidity"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
}

func FetchWithStatus(city string) (*model.Weather, int, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		return nil, http.StatusBadGateway, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("failed to close response body: %v", cerr)
		}
	}()

	switch resp.StatusCode {
	case 400:
		return nil, http.StatusBadRequest, fmt.Errorf("Invalid city name")
	case 404:
		return nil, http.StatusNotFound, fmt.Errorf("City not found")
	case 200:
		// OK
	default:
		return nil, http.StatusBadGateway, fmt.Errorf("Weather API returned unexpected status")
	}

	var data weatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("Failed to parse weather data")
	}

	result := &model.Weather{
		Temperature: data.Current.TempC,
		Humidity:    data.Current.Humidity,
		Description: data.Current.Condition.Text,
	}

	return result, http.StatusOK, nil
}
