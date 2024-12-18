package models

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	gorm.Model
	BaseValidator
	UserID    uint
	Type      string `gorm:"not null"`
	Message   string `gorm:"not null"`
	IsRead    bool   `gorm:"default:false"`
	CreatedAt time.Time
}

func (q *Notification) Validate() error {
	return q.BaseValidator.Validate(q)
}
