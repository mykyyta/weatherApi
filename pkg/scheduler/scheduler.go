package scheduler

import (
	"log"
	"time"

	"weatherApi/internal/model"
	"weatherApi/pkg/email"
	"weatherApi/pkg/weatherapi"

	"gorm.io/gorm"
)

var DB *gorm.DB

func SetDB(db *gorm.DB) {
	DB = db
}

func StartWeatherScheduler() {
	log.Println("[Scheduler] started")

	now := time.Now()
	sleepDuration := time.Until(now.Truncate(time.Minute).Add(time.Minute))
	log.Printf("[Scheduler] waiting %v to align to :00 seconds\n", sleepDuration)
	time.Sleep(sleepDuration)

	ticker := time.NewTicker(1 * time.Minute)
	for {
		now := time.Now()
		sendDaily := now.Minute()%5 == 0 // або: now.Hour() == 9

		log.Println("[Scheduler] running tick", now.Format("15:04:05"))
		sendWeatherUpdates("hourly")
		if sendDaily {
			sendWeatherUpdates("daily")
		}

		<-ticker.C
	}
}
func sendWeatherUpdates(frequency string) {
	var subs []model.Subscription
	if err := DB.Where("is_confirmed = ? AND is_unsubscribed = ? AND frequency = ?", true, false, frequency).Find(&subs).Error; err != nil {
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

func ProcessSubscription(sub model.Subscription) error {
	weather, _, err := weatherapi.FetchWithStatus(sub.City)
	if err != nil {
		return err
	}
	return email.SendWeatherEmail(sub.Email, weather, sub.City, sub.Token)
}
