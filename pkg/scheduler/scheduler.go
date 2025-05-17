package scheduler

import (
	"log"
	"time"

	"weatherApi/internal/model"
	"weatherApi/pkg/email"
	"weatherApi/pkg/weatherapi"

	"gorm.io/gorm"
)

// DB is the scheduler's database instance.
// Must be set via SetDB() before StartWeatherScheduler is called.
var DB *gorm.DB

// SetDB assigns a GORM database instance to the scheduler.
// This allows decoupling from the main DB package for testability or modularity.
func SetDB(db *gorm.DB) {
	DB = db
}

// StartWeatherScheduler starts the periodic task that sends weather updates
// according to subscription frequency.
// It aligns execution to :00 seconds of each minute to maintain precision.
func StartWeatherScheduler() {
	log.Println("[Scheduler] started")

	// Wait until the next full minute (e.g., xx:01:00) to align the ticker
	now := time.Now()
	sleepDuration := time.Until(now.Truncate(time.Minute).Add(time.Minute))
	log.Printf("[Scheduler] waiting %v to align to :00 seconds\n", sleepDuration)
	time.Sleep(sleepDuration)

	ticker := time.NewTicker(1 * time.Minute)
	for {
		now := time.Now()

		// For demo/testing: send "daily" every 5 minutes
		// In production: use e.g. now.Hour() == 9
		sendDaily := now.Minute()%5 == 0

		log.Println("[Scheduler] running tick", now.Format("15:04:05"))
		sendWeatherUpdates("hourly")
		if sendDaily {
			sendWeatherUpdates("daily")
		}

		<-ticker.C
	}
}

// sendWeatherUpdates fetches all active subscriptions with the given frequency
// and sends weather updates for each one via email.
func sendWeatherUpdates(frequency string) {
	var subs []model.Subscription

	if err := DB.Where(
		"is_confirmed = ? AND is_unsubscribed = ? AND frequency = ?",
		true, false, frequency,
	).Find(&subs).Error; err != nil {
		log.Printf("[Scheduler] Failed to query subscriptions: %v", err)
		return
	}

	for _, sub := range subs {
		if err := ProcessSubscription(sub); err != nil {
			log.Printf("[Scheduler] Failed to process %s: %v", sub.Email, err)
		} else {
			log.Printf("[Scheduler] Weather sent to %s", sub.Email)
		}
	}
}

// ProcessSubscription fetches the weather for a single subscription
// and sends the email using the stored unsubscribe token.
func ProcessSubscription(sub model.Subscription) error {
	weather, _, err := weatherapi.FetchWithStatus(sub.City)
	if err != nil {
		return err
	}
	return email.SendWeatherEmail(sub.Email, weather, sub.City, sub.Token)
}
