package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	BaseValidator
	Content       string `gorm:"not null"`
	ChatID        uint   `gorm:"not null"`
	UserID        uint   `gorm:"not null"`
	AttachmentURL *string
	User          User `gorm:"foreignKey:UserID"`
}

func (q *Message) Validate() error {
	return q.BaseValidator.Validate(q)
}
