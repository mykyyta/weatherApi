package scheduler

import (
	"log"
	"time"
	"weatherApi/internal/db"
	"weatherApi/internal/model"
	"weatherApi/pkg/email"
	"weatherApi/pkg/weatherapi"
)

func StartWeatherScheduler() {
	log.Println("[Scheduler] started")
	ticker := time.NewTicker(1 * time.Minute) // ticker := time.NewTicker(1 * time.Hour)
	for {
		now := time.Now()
		sendDaily := now.Minute()%5 == 0 // sendDaily := now.Hour() == 9

		log.Println("[Scheduler] running tick", now.Format("15:04"))
		sendWeatherUpdates("hourly")
		if sendDaily {
			sendWeatherUpdates("daily")
		}

		<-ticker.C
	}
}

func sendWeatherUpdates(frequency string) {
	var subs []model.Subscription
	if err := db.DB.Where("is_confirmed = ? AND is_unsubscribed = ? AND frequency = ?", true, false, frequency).Find(&subs).Error; err != nil {
		log.Printf("[Scheduler] Failed to query subscriptions: %v", err)
		return
	}

	for _, sub := range subs {
		weather, status, err := weatherapi.FetchWithStatus(sub.City)
		if err != nil {
			log.Printf("[Scheduler] Weather fetch error for %s: %v (status: %d)", sub.City, err, status)
			continue
		}

		if err := email.SendWeatherEmail(sub.Email, weather, sub.City); err != nil {
			log.Printf("[Scheduler] Failed to send email to %s: %v", sub.Email, err)
		} else {
			log.Printf("[Scheduler] Weather sent to %s", sub.Email)
		}
	}
}
