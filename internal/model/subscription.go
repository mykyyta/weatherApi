package model

import (
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email       string    `gorm:"not null;uniqueIndex"`
	City        string    `gorm:"not null"`
	Frequency   string    `gorm:"type:text;check:frequency IN ('daily','hourly');not null"`
	IsConfirmed bool      `gorm:"default:false"`
	Token       string    `gorm:"not null;uniqueIndex"`
	CreatedAt   time.Time
}
