package email

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"weatherApi/internal/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(toEmail, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail("weatherApp", os.Getenv("EMAIL_FROM"))
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
	subject := "Підтвердіть вашу підписку на погодні сповіщення"

	confirmURL := fmt.Sprintf("http://localhost:8080/api/confirm/%s", token)
	plainText := "Будь ласка, підтвердіть вашу підписку: " + confirmURL
	htmlContent := fmt.Sprintf(`<p>Натисніть нижче для підтвердження вашої підписки:</p><p><a href="%s">Підтвердити підписку</a></p>`, confirmURL)

	return SendEmail(toEmail, subject, plainText, htmlContent)
}

func SendWeatherEmail(toEmail string, weather *model.Weather, city string, token string) error {
	caser := cases.Title(language.English)
	subject := fmt.Sprintf("Ваше оновлення погоди для %s", caser.String(city))

	unsubscribeURL := fmt.Sprintf("http://localhost:8080/api/unsubscribe/%s", token)

	plainText := fmt.Sprintf(
		"Вітаємо!\n\nПоточна погода в %s:\nТемпература: %.1f°C\nВологість: %d%%\nОпис: %s\n\nЯкщо бажаєте скасувати підписку, перейдіть за посиланням: %s",
		caser.String(city), weather.Temperature, weather.Humidity, weather.Description, unsubscribeURL,
	)

	htmlContent := fmt.Sprintf(
		`<h2>Погода в %s</h2>
		<p><strong>Температура:</strong> %.1f°C</p>
		<p><strong>Вологість:</strong> %d%%</p>
		<p><strong>Опис:</strong> %s</p>
		<hr>
		<p style="font-size:small">Не хочете більше отримувати? <a href="%s">Відписатися</a></p>`,
		caser.String(city), weather.Temperature, weather.Humidity, weather.Description, unsubscribeURL,
	)

	return SendEmail(toEmail, subject, plainText, htmlContent)
}
