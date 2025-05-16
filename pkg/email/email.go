package email

import (
	"fmt"
	"os"

	"weatherApi/internal/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(toEmail, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail("WeatherBot", os.Getenv("EMAIL_FROM"))
	to := mail.NewEmail("User", toEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid failed with status %d: %s", response.StatusCode, response.Body)
	}

	return nil
}

func SendConfirmationEmail(toEmail, token string) error {
	subject := "Confirm your Weather Subscription"

	confirmURL := fmt.Sprintf("http://localhost:8080/api/confirm/%s", token)
	plainText := "Please confirm your subscription: " + confirmURL
	htmlContent := fmt.Sprintf(`<p>Click below to confirm your subscription:</p><p><a href="%s">Confirm Subscription</a></p>`, confirmURL)

	return SendEmail(toEmail, subject, plainText, htmlContent)
}

func SendWeatherEmail(toEmail string, weather *model.Weather, city string) error {
	subject := fmt.Sprintf("Your Weather Update for %s", city)

	plainText := fmt.Sprintf(
		"Hello!\n\nCurrent weather in %s:\nTemperature: %.1f°C\nHumidity: %d%%\nDescription: %s\n",
		city, weather.Temperature, weather.Humidity, weather.Description,
	)

	htmlContent := fmt.Sprintf(
		`<h2>Weather in %s</h2><p><strong>Temperature:</strong> %.1f°C</p><p><strong>Humidity:</strong> %d%%</p><p><strong>Description:</strong> %s</p>`,
		city, weather.Temperature, weather.Humidity, weather.Description,
	)

	return SendEmail(toEmail, subject, plainText, htmlContent)
}
