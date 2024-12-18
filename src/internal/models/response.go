package models

import (
	"gorm.io/gorm"
	"time"
)

type Response struct {
	gorm.Model
	BaseValidator
	UserID          uint
	QuestionnaireID uint
	QuestionID      uint
	OptionID        uint
	Answer          string `gorm:"not null"`
	IsFinalized     bool   `gorm:"default: false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (q *Response) Validate() error {
	return q.BaseValidator.Validate(q)
}
