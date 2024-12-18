package models

import (
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	BaseValidator
	Content         string `gorm:"not null"`
	SenderUserID    uint   `gorm:"not null"`
	ReceiverUserID  uint   `gorm:"not null"`
	QuestionnaireID uint   `gorm:"not null"`
	Type            string `gorm:"not null"` // "private" or "group"
	Status          string `gorm:"not null"` // "active" or "inactive"
	Sender          User   `gorm:"foreignKey:SenderUserID"`
	Receiver        User   `gorm:"foreignKey:ReceiverUserID"`
}

func (q *Chat) Validate() error {
	return q.BaseValidator.Validate(q)
}
