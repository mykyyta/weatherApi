package model

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email          string    `gorm:"not null;uniqueIndex" json:"email"`
	City           string    `gorm:"not null" json:"city"`
	Frequency      string    `gorm:"type:text;check:frequency IN ('daily','hourly');not null" json:"frequency"`
	IsConfirmed    bool      `gorm:"default:false" json:"is_confirmed"`
	IsUnsubscribed bool      `gorm:"default:false" json:"is_unsubscribed"`
	Token          string    `gorm:"not null" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}
