package models

import (
	"time"

	"gorm.io/gorm"
)

// Email struct.
type Email struct {
	gorm.Model
	To      string
	Subject string
	Message string
	SendAt  time.Time
	IsSent  bool
}
